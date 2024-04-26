package gstrlog

import (
	"errors"

	log "github.com/omgolab/go-commons/pkg/log"
	filter "github.com/omgolab/go-commons/pkg/log/custom"
)

type StringLogger interface {
	filter.FilterLogger
	GetStringLogs() []string
	AppendString(str string)
	AppendStringErr(str string, err error)
}

type arrayCfg struct {
	filter.FilterLogger
	data []string
}

func (c *arrayCfg) Write(b []byte) (n int, err error) {
	c.data = append(c.data, string(b))
	return len(b), nil
}

func (c *arrayCfg) GetStringLogs() []string {
	return c.data
}

func (c *arrayCfg) AppendString(str string) {
	c.TagLog(str, log.DebugLevel, nil, 3)
}

func (c *arrayCfg) AppendStringErr(str string, err error) {
	c.TagLog(str, log.ErrorLevel, err, 3)
}

// New creates a new instance of the StringLogger interface, which is a logger implementation that logs
// messages to an array and allows retrieving the logged messages.
func New(filterOpts []filter.FilterOption, logOpts ...log.LogOption) (StringLogger, error) {
	// Create a new instance of the arrayCfg struct
	sl := &arrayCfg{}

	// Create a new filter logger with the provided options
	var err error
	sl.FilterLogger, err = filter.New("string-log", sl, filterOpts, logOpts...)
	if err != nil {
		return nil, err
	}

	// Check if the filter logger is nil
	if sl.FilterLogger == nil {
		return nil, errors.New("filter.New returned a nil FilterLogger")
	}

	return sl, nil
}
