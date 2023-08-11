package log_utils

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	csv_utils "github.com/omar391/go-commons/pkg/csv"
	fu "github.com/omar391/go-commons/pkg/file"
	stream_utils "github.com/omar391/go-commons/pkg/stream"
	"github.com/rs/zerolog"
)

const CsvIndicator = "~csv~"

type CsvLogger struct {
	cw                       *csvWriter
	logger                   zerolog.Logger
	truncateOnHeadersMissing bool
	disableTimestamp         bool
	enableLevel             bool
	enableCaller             bool
	enableError              bool
}

type csvWriter struct {
	// default is ','
	comma   string
	csvFile *os.File
}

// this writer validates the console output and if ok then push to file
func (cw *csvWriter) Write(b []byte) (n int, err error) {
	ln := len(b) // used for returning the original length
	bi := []byte(CsvIndicator)
	l := len(bi) + 1 // 1 for the comma or a ending newline

	for {
		m := bytes.LastIndex(b, bi)
		if m < 0 {
			break
		}

		j := m + l
		if j < len(b) && b[j] == ' ' {
			j++
		} else if j == len(b) {
			j--
		}

		copy(b[m:], b[j:])
		b = b[:len(b)-(j-m)]
	}

	bl := len(b)
	if bl < ln {
		if bl > 2 && b[bl-3] == cw.comma[0] && b[bl-2] == ' ' {
			b[bl-3] = '\n'
			b = b[:bl-2]
		}
		_, err = cw.csvFile.Write(b)
	}

	return ln, err
}

func (c *CsvLogger) FileName() string {
	return c.cw.csvFile.Name()
}

func (c *CsvLogger) Csv(str ...string) {
	// this is a special string to indicate csv write
	c.logger.Info().Msg(strings.Join(str, CsvIndicator))
}

func (c *CsvLogger) Logger() *zerolog.Logger {
	return &c.logger
}

func (c *CsvLogger) Close() {
	c.cw.csvFile.Close()
}

// NewCsvLogger creates a dual logger with csv and console output
// remember to close the file after use
func NewCsvLogger(csvPath string, headers []string, opts ...LogOption) (*CsvLogger, error) {
	// create a csv c
	c := &CsvLogger{
		cw: &csvWriter{
			comma: ",",
		},
	}

	// update the options
	var err error
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return c, err
		}
	}
	// set a default csv path
	c.cw.csvFile, err = c.getCsvFile(csvPath, headers)
	if err != nil {
		return c, err
	}

	// Setup the csv logger
	cl := zerolog.ConsoleWriter{
		Out:     c.cw,
		NoColor: true,
		FormatTimestamp: func(i interface{}) string {
			if c.disableTimestamp {
				return ""
			}
			format := "2006-01-02 05:04PM"
			v := i.(string)
			ts, err := time.ParseInLocation(format, v, time.Local)
			if err != nil {
				return v + c.cw.comma
			} else {
				return ts.Local().Format(time.DateTime) + c.cw.comma
			}
		},
		FormatLevel: func(i interface{}) string {
			if !c.enableLevel {
				return ""
			}
			return strings.ToUpper(i.(string))[0:3] + c.cw.comma
		},
		FormatMessage: func(i interface{}) string {
			s, b := i.(string)
			if i == nil || !b {
				return ""
			}

			// split the csv indicator and return if CsvIndicator's not found
			as := strings.Split(s, CsvIndicator)
			if len(as) < 2 {
				return ""
			}

			// remove the csv indicator and append a comma and escape the csv if needed
			// last CsvIndicator is needed for the final file writer
			return stream_utils.New(as).Reduce("", func(o, n string) string {
				if len(o) == 0 {
					return string(csv_utils.Escape([]byte(n)))
				}
				if len(n) == 0 {
					return o
				}
				return o + c.cw.comma + " " + string(csv_utils.Escape([]byte(n)))
			}) + CsvIndicator
		},
		FormatCaller: func(i interface{}) string {
			if !c.enableCaller {
				return ""
			}
			var s string
			if cc, ok := i.(string); ok {
				s = cc
			}
			if len(s) > 0 {
				if cwd, err := os.Getwd(); err == nil {
					if rel, err := filepath.Rel(cwd, s); err == nil {
						s = rel
					}
				}
				s += " >"
			}
			return s + c.cw.comma
		},
		FormatFieldName: func(i interface{}) string {
			// don't allow any extra fields other than the ones specified in the headers
			return ""
		},
		FormatFieldValue: func(i interface{}) string {
			// don't allow any extra fields other than the ones specified in the headers
			return ""
		},
		FormatErrFieldName: func(i interface{}) string {
			// no need to write "error="
			return ""
		},
		FormatErrFieldValue: func(i interface{}) string {
			if !c.enableError {
				return ""
			}
			return fmt.Sprintf("%s%s ", i, c.cw.comma)
		},
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}

	ctx := zerolog.New(zerolog.MultiLevelWriter(cl, zerolog.NewConsoleWriter())).With()
	if !c.disableTimestamp {
		ctx = ctx.Timestamp()
	}
	if c.enableCaller {
		ctx = ctx.Caller()
	}
	c.logger = ctx.Logger()

	return c, nil
}

func (c *CsvLogger) prepareCsvHeaderLine(headers []string) string {
	preHeaders := []string{}
	if !c.disableTimestamp {
		preHeaders = append(preHeaders, "Timestamp")
	}
	if c.enableLevel {
		preHeaders = append(preHeaders, "Level")
	}
	if c.enableCaller {
		preHeaders = append(preHeaders, "Caller")
	}

	// append message headers
	preHeaders = append(preHeaders, headers...)

	// finally append error headers
	// since error are appended at the end by console writer
	if c.enableError {
		preHeaders = append(preHeaders, "Error")
	}
	return strings.Join(preHeaders, c.cw.comma+" ") + "\n"
}

// create a csv file writer for the hook
// remember to call: file.Close()
func (c *CsvLogger) getCsvFile(path string, headers []string) (*os.File, error) {
	opts := []fu.OpenOption{
		// create a new file with incremental _number suffix
		fu.WithIncrementalSuffixIfExists(func(path string, fi fs.FileInfo) bool {
			// false: don't increment the file name
			// true: increment the file name

			// if the file doesn't exists
			if fi == nil {
				return false
			}

			b, err := fu.ContainsAllTexts(path, 1, 1, headers...)
			if err != nil {
				return false
			}

			// if the headers are not found and the truncate option is set
			if !b && c.truncateOnHeadersMissing {
				return false
			}

			// if the file is empty
			if fi.Size() == 0 {
				return false
			}

			return !b
		}),
	}
	if c.truncateOnHeadersMissing {
		opts = append(opts, fu.WithTruncate())
	}

	// open the file
	f, err := fu.OpenFile(path, opts...)
	if err != nil {
		return nil, err
	}

	// finally, write the headers if the file is empty
	s, _ := f.Stat()
	if s.Size() == 0 {
		_, err = f.WriteString(c.prepareCsvHeaderLine(headers))
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func WithDelimiter(delim rune) LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.cw.comma = string(delim)
		return nil
	}
}

func WithTruncateOnHeadersMissing() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.truncateOnHeadersMissing = true
		return nil
	}
}

func WithDisableTimestamp() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.disableTimestamp = true
		return nil
	}
}

func WithEnableLevel() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableLevel = true
		return nil
	}
}

func WithEnableCaller() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableCaller = true
		return nil
	}
}

func WithEnableError() LogOption {
	return func(h any) error {
		ch := h.(*CsvLogger)
		ch.enableError = true
		return nil
	}
}
