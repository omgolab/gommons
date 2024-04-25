package gclog

import (
	"io"
	"os"

	gcfile "github.com/omgolab/go-commons/pkg/file/open"
)

type LogOption func(*logCfg) error

// Option setters:
// Note: we use these option when
// - we need to replace the default logger
// - we usually don't need to update those attrs dynamically later

func WithFileLogger(filename string) LogOption {
	return func(l *logCfg) error {
		f, err := gcfile.OpenFile(filename)
		if err != nil {
			return nil
		}
		l.sc.writers = append(l.sc.writers, f)
		return nil
	}
}

func WithJsonStdOut() LogOption {
	return func(l *logCfg) error {
		l.sc.writers[0] = os.Stdout // update zero index to stdout
		return nil
	}
}

func WithMultiLogger(ws ...io.Writer) LogOption {
	return func(l *logCfg) error {
		l.sc.writers = append(l.sc.writers, ws...)
		return nil
	}
}

func WithTimestampFormat(format string) LogOption {
	return func(l *logCfg) error {
		l.sc.timeFormat = format
		return nil
	}
}

func WithDefaultLogLevel(level LogLevel) LogOption {
	return func(l *logCfg) error {
		l.sc.minLogLevel = level
		return nil
	}
}
