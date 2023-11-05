package gccsvlog

import (
	"io/fs"
	"os"
	"strings"

	file "github.com/omar391/go-commons/pkg/file"
	fo "github.com/omar391/go-commons/pkg/file/open"
	log "github.com/omar391/go-commons/pkg/log"
	tl "github.com/omar391/go-commons/pkg/log/gccustlog"
)

type CsvLogger interface {
	tl.TaggedLogger
	FileName() string
	Csv(str ...string)
	CsvErr(err error, str ...string)
	NewFork(fileName string) (CsvLogger, error)
}

type csvCfg struct {
	tl.TaggedLogger
	file                     *os.File
	truncateOnHeadersMissing bool
	tlOpts                   []tl.TaggedLoggerOption
	headers                  []string
}

func (c *csvCfg) NewFork(fileName string) (CsvLogger, error) {
	nc := *c
	cl := &nc
	err := cl.prepareCsvFile(fileName)
	if err != nil {
		return nil, err
	}
	err = cl.writeCsvHeaders()

	return cl, err
}

func (c *csvCfg) FileName() string {
	return c.file.Name()
}

func (c *csvCfg) Csv(str ...string) {
	// this is a special string to indicate csv write
	c.LogTag(strings.Join(str, c.GetDelimiter()), log.DebugLevel, nil, 3)
}

func (c *csvCfg) CsvErr(err error, str ...string) {
	// this is a special string to indicate csv write
	c.LogTag(strings.Join(str, c.GetDelimiter()), log.ErrorLevel, err, 3)
}

// New creates a dual logger with csv and console output
func New(csvPath string, opts ...CsvOption) (CsvLogger, error) {
	// create a csv c
	c := &csvCfg{}

	// update the options
	var err error
	for _, opt := range opts {
		opt(c)
	}

	// set a default csv path
	err = c.prepareCsvFile(csvPath)
	if err != nil {
		return nil, err
	}

	// update the tagged logger
	c.TaggedLogger, err = tl.New("csv", c.file, c.tlOpts...)
	if err != nil {
		return nil, err
	}

	// write headers
	err = c.writeCsvHeaders()

	return c, err
}

// create a csv file writer for the hook
// remember to call: file.Close()
func (c *csvCfg) prepareCsvFile(path string) error {
	opts := []fo.OpenOption{
		// create a new file with incremental _number suffix
		fo.WithIncrementalSuffixIfExists(func(path string, fi fs.FileInfo) bool {
			// false: don't increment the file name
			// true: increment the file name

			// if the file doesn't exists
			if fi == nil {
				return false
			}

			b, err := file.ContainsAllTexts(path, 1, 1, c.headers...)
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
	var err error
	c.file, err = fo.OpenFile(path, opts...)
	return err
}

func (c *csvCfg) writeCsvHeaders() (err error) {
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
	preHeaders = append(preHeaders, c.headers...)

	// finally append error headers
	// since error are appended at the end by console writer
	if c.IsErrorFormatterEnabled() {
		preHeaders = append(preHeaders, "Error")
	}
	h := strings.Join(preHeaders, c.GetDelimiter()+" ") + "\n"

	_, err = c.file.WriteString(h)
	return err
}
