package glog

import (
	"io"
	"os"

	gfile "github.com/omgolab/go-commons/pkg/file/open"
	"github.com/rs/zerolog"
)

type LogOption func(*logCfg) error

// Option setters:
// Note: we use these option when
// - we need to replace the default logger
// - we usually don't need to update those attrs dynamically later

func WithFileLogger(filename string) LogOption {
	return func(l *logCfg) error {
		f, err := gfile.OpenFile(filename)
		if err != nil {
			return nil
		}
		l.sc.writers = append(l.sc.writers, f)
		return nil
	}
}

func WithJsonStdOut() LogOption {
	return func(l *logCfg) error {
		if len(l.sc.writers) == 0 {
			l.sc.writers = []io.Writer{os.Stdout}
		} else {
			l.sc.writers[0] = os.Stdout
		}
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
		for i, w := range l.sc.writers {
			if cw, ok := w.(*zerolog.ConsoleWriter); ok {
				cw.TimeFormat = format
				l.sc.writers[i] = cw
			}
		}
		return nil
	}
}

func WithDefaultLogLevel(level LogLevel) LogOption {
	return func(l *logCfg) error {
		l.sc.minLogLevel = level
		return nil
	}
}
