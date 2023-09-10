package gclog

import (
	"io"
	"os"

	gccollections "github.com/omar391/go-commons/pkg/collections"
	file "github.com/omar391/go-commons/pkg/file/open"
	"github.com/rs/zerolog"
)

type LogFields map[string]interface{}
type LogLevel int

const (
	NoLevel LogLevel = iota
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

var logToZerologMap = map[LogLevel]zerolog.Level{
	NoLevel:    zerolog.NoLevel,
	TraceLevel: zerolog.TraceLevel,
	DebugLevel: zerolog.DebugLevel,
	InfoLevel:  zerolog.InfoLevel,
	WarnLevel:  zerolog.WarnLevel,
	ErrorLevel: zerolog.ErrorLevel,
	FatalLevel: zerolog.FatalLevel,
	PanicLevel: zerolog.PanicLevel,
}

type Logger interface {
	SetMinLevel(level LogLevel) Logger
	SetContext(keyword string) Logger
	Disable() Logger
	Println(msg ...any) Logger
	Trace(msg string, fields ...LogFields) Logger
	Debug(msg string, fields ...LogFields) Logger
	Info(msg string, fields ...LogFields) Logger
	Warn(msg string, fields ...LogFields) Logger
	Error(msg string, err error, fields ...LogFields) Logger
	Fatal(msg string, err error, fields ...LogFields) Logger
	Panic(msg string, err error, fields ...LogFields) Logger
}

type OptionSetter func(*logger) error

type logger struct {
	logger       zerolog.Logger
	minLevel     LogLevel
	contextAlias string
	disabled     bool
	timestampOff bool
}

func NewLogger(options ...OptionSetter) (Logger, error) {
	l := &logger{
		logger:   zerolog.New(os.Stdout).With().Timestamp().Logger(),
		minLevel: DebugLevel,
	}

	for _, opt := range options {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// Option setters:
// Note: we use these option when
// - we need to replace the default logger
// - we usually don't need to update those attrs dynamically later

func WithFileLogger(filename string, jsonStdOut bool) OptionSetter {
	return func(l *logger) error {
		f, err := file.OpenFile(filename)
		if err != nil {
			return nil
		}

		var cl io.Writer
		cl = zerolog.NewConsoleWriter()
		if jsonStdOut {
			cl = os.Stdout
		}

		lg := zerolog.New(zerolog.MultiLevelWriter(cl, f))
		if !l.timestampOff {
			lg = lg.With().Timestamp().Logger()
		}

		l.logger = lg
		return nil
	}
}

func WithoutTimestamp() OptionSetter {
	return func(l *logger) error {
		l.logger = l.logger.With().Logger()
		l.timestampOff = true
		return nil
	}
}

func WithDefaultLevel(level LogLevel) OptionSetter {
	return func(l *logger) error {
		l.minLevel = level
		return nil
	}
}

func WithDefaultContext(keyword string) OptionSetter {
	return func(l *logger) error {
		l.contextAlias = keyword
		return nil
	}
}

// Logger methods:

func (l *logger) SetMinLevel(level LogLevel) Logger {
	l.minLevel = level
	return l
}

func (l *logger) SetContext(keyword string) Logger {
	l.contextAlias = keyword
	return l
}

func (l *logger) logEvent(level LogLevel, msg string, err error, fields []LogFields) Logger {
	if l.disabled || level < l.minLevel {
		return l
	}

	event := l.logger.WithLevel(logToZerologMap[level])

	if len(fields) > 0 {
		lfs := gccollections.MergeMaps[LogFields](fields...)
		if l.contextAlias != "" {
			lfs["context"] = l.contextAlias
		}
		if lfs != nil {
			event = event.Fields(lfs)
		}
	}

	if level >= WarnLevel {
		event = event.Stack()
	}

	if err != nil {
		event = event.Err(err)
	}

	event.Msg(msg)
	return l
}

func (l *logger) Trace(msg string, fields ...LogFields) Logger {
	return l.logEvent(TraceLevel, msg, nil, fields)
}

func (l *logger) Debug(msg string, fields ...LogFields) Logger {
	return l.logEvent(DebugLevel, msg, nil, fields)
}

func (l *logger) Info(msg string, fields ...LogFields) Logger {
	return l.logEvent(InfoLevel, msg, nil, fields)
}

func (l *logger) Warn(msg string, fields ...LogFields) Logger {
	return l.logEvent(WarnLevel, msg, nil, fields)
}

func (l *logger) Error(msg string, err error, fields ...LogFields) Logger {
	return l.logEvent(ErrorLevel, msg, err, fields)
}

func (l *logger) Fatal(msg string, err error, fields ...LogFields) Logger {
	return l.logEvent(FatalLevel, msg, err, fields)
}

func (l *logger) Panic(msg string, err error, fields ...LogFields) Logger {
	return l.logEvent(PanicLevel, msg, err, fields)
}

func (l *logger) Disable() Logger {
	l.disabled = true
	return l
}

func (l *logger) Println(msg ...any) Logger {
	if !l.disabled {
		l.logger.Print(msg...)
	}
	return l
}
