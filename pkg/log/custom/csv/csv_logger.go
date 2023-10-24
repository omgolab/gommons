package gccsvlog

import (
	"io/fs"
	"os"
	"strings"

	file "github.com/omar391/go-commons/pkg/file"
	fo "github.com/omar391/go-commons/pkg/file/open"
	log "github.com/omar391/go-commons/pkg/log"
	filter "github.com/omar391/go-commons/pkg/log/custom"
)

type CsvLogger interface {
	filter.FilterLogger
	FileName() string
	Csv(str ...string)
	CsvErr(err error, str ...string)
}

type csvCfg struct {
	filter.FilterLogger
	file                     *os.File
	truncateOnHeadersMissing bool
}

func (c *csvCfg) FileName() string {
	return c.file.Name()
}

func (c *csvCfg) Csv(str ...string) {
	// this is a special string to indicate csv write
	c.TagLog(strings.Join(str, c.GetDelimiter()), log.DebugLevel, nil, 3)
}

func (c *csvCfg) CsvErr(err error, str ...string) {
	// this is a special string to indicate csv write
	c.TagLog(strings.Join(str, c.GetDelimiter()), log.ErrorLevel, err, 3)
}

// New creates a dual logger with csv and console output
func New(csvPath string, headers []string, opts []CsvOption, filterOpts []filter.FilterOption, logOpts ...log.LogOption) (CsvLogger, error) {
	// create a csv c
	c := &csvCfg{}

	// update the options
	var err error
	for _, opt := range opts {
		err = opt(c)
		if err != nil {
			return c, err
		}
	}

	// set a default csv path
	c.file, err = c.getCsvFile(csvPath, headers)
	if err != nil {
		return nil, err
	}

	// update the filter logger
	c.FilterLogger, err = filter.New("csv", c.file, filterOpts, logOpts...)
	if err != nil {
		return nil, err
	}

	// write headers
	err = c.writeCsvHeaders(headers)

	return c, err
}

func (c *csvCfg) writeCsvHeaders(headers []string) (err error) {
	// finally, write the headers if the file is empty
	s, _ := c.file.Stat()
	if s.Size() != 0 {
		return nil
	}

	preHeaders := []string{}
	if c.IsTimestampFormatterEnabled() {
		preHeaders = append(preHeaders, "Timestamp")
	}
	if c.IsLevelFormatterEnabled() {
		preHeaders = append(preHeaders, "Level")
	}
	if c.IsCallerFormatterEnabled() {
		preHeaders = append(preHeaders, "Caller")
	}

	// append message headers
	preHeaders = append(preHeaders, headers...)

	// finally append error headers
	// since error are appended at the end by console writer
	if c.IsErrorFormatterEnabled() {
		preHeaders = append(preHeaders, "Error")
	}
	h := strings.Join(preHeaders, c.GetDelimiter()+" ") + "\n"

	_, err = c.file.WriteString(h)
	return err
}

// create a csv file writer for the hook
// remember to call: file.Close()
func (c *csvCfg) getCsvFile(path string, headers []string) (*os.File, error) {
	opts := []fo.OpenOption{
		// create a new file with incremental _number suffix
		fo.WithIncrementalSuffixIfExists(func(path string, fi fs.FileInfo) bool {
			// false: don't increment the file name
			// true: increment the file name

			// if the file doesn't exists
			if fi == nil {
				return false
			}

			b, err := file.ContainsAllTexts(path, 1, 1, headers...)
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
		opts = append(opts, fo.WithTruncate())
	}

	// open the file
	f, err := fo.OpenFile(path, opts...)
	if err != nil {
		return nil, err
	}

	return f, nil
}
