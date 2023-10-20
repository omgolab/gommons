package gclog

import (
	"io"
	"os"

	gcfile "github.com/omar391/go-commons/pkg/file/open"
)

type OptionSetter func(*logger) error

// Option setters:
// Note: we use these option when
// - we need to replace the default logger
// - we usually don't need to update those attrs dynamically later

func WithFileLogger(filename string) OptionSetter {
	return func(l *logger) error {
		f, err := gcfile.OpenFile(filename)
		if err != nil {
			return nil
		}
		l.sc.writers = append(l.sc.writers, f)
		return nil
	}
}

func WitJsonStdOut() OptionSetter {
	return func(l *logger) error {
		l.sc.writers[0] = os.Stdout // update zero index to stdout
		return nil
	}
}

func WithMultiLogger(ws ...io.Writer) OptionSetter {
	return func(l *logger) error {
		l.sc.writers = append(l.sc.writers, ws...)
		return nil
	}
}

func WithTimestampFormat(format string) OptionSetter {
	return func(l *logger) error {
		l.sc.timeFormat = format
		return nil
	}
}

func WithDefaultLogLevel(level LogLevel) OptionSetter {
	return func(l *logger) error {
		l.sc.minLogLevel = level
		return nil
	}
}
