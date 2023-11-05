package gcstrlog

import (
	"errors"

	log "github.com/omar391/go-commons/pkg/log"
	tl "github.com/omar391/go-commons/pkg/log/gccustlog"
)

type StringLogger interface {
	tl.TaggedLogger
	GetStringLogs() []string
	AppendString(str string)
	AppendStringErr(str string, err error)
	NewFork() StringLogger
}

type arrayCfg struct {
	tl.TaggedLogger
	data            []string
	initialDataSize int
	tlOpts          []tl.TaggedLoggerOption
}

func (c *arrayCfg) Write(b []byte) (n int, err error) {
	c.data = append(c.data, string(b))
	return len(b), nil
}

func (c *arrayCfg) GetStringLogs() []string {
	return c.data
}

func (c *arrayCfg) AppendString(str string) {
	c.LogTag(str, log.DebugLevel, nil, 3)
}

func (c *arrayCfg) AppendStringErr(str string, err error) {
	c.LogTag(str, log.ErrorLevel, err, 3)
}

// create a new forked logger
func (c *arrayCfg) NewFork() StringLogger {
	nc := *c
	n := &nc
	n.data = make([]string, c.initialDataSize)
	return n
}

// New creates a new instance of the StringLogger interface, which is a logger implementation that logs
// messages to an array and allows retrieving the logged messages.
func New(opts ...StringLoggerOption) (StringLogger, error) {
	// Create a new instance of the arrayCfg struct
	sl := &arrayCfg{
		initialDataSize: 10,
	}

	// apply options
	for _, opt := range opts {
		opt(sl)
	}

	// Create a new tagged logger with the provided options
	var err error
	sl.TaggedLogger, err = tl.New("string-log", sl, sl.tlOpts...)
	if err != nil {
		return nil, err
	}

	// Check if the tagged logger is nil
	if sl.TaggedLogger == nil {
		return nil, errors.New("tagged logger - New, returned a nil TaggedLogger")
	}

	return sl, nil
}
