package gclog

import (
	"context"
	"io"
	"os"

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
	SetContext(ctx context.Context) Logger
	SetKeyword(keyword string) Logger
	Trace(msg string, fields LogFields) Logger
	Debug(msg string, fields LogFields) Logger
	Info(msg string, fields LogFields) Logger
	Warn(msg string, fields LogFields) Logger
	Error(msg string, fields LogFields) Logger
	Fatal(msg string, fields LogFields) Logger
	Panic(msg string, fields LogFields) Logger
	Disable() Logger
	Println(msg string) Logger
}

type OptionSetter func(*logger)

type logger struct {
	logger       zerolog.Logger
	minLevel     LogLevel
	contextAlias string
	disabled     bool
}

func New(level LogLevel, options ...OptionSetter) Logger {
	l := &logger{
		logger:   zerolog.New(os.Stderr).With().Timestamp().Logger(),
		minLevel: level,
	}

	for _, opt := range options {
		opt(l)
	}

	return l
}

// Option setters:

func WithOutput(output io.Writer) OptionSetter {
	return func(l *logger) {
		l.logger = l.logger.Output(output)
	}
}

func WithoutTimestamp() OptionSetter {
	return func(l *logger) {
		l.logger = l.logger.With().Logger()
	}
}

func WithDefaultLevel(level LogLevel) OptionSetter {
	return func(l *logger) {
		zerolog.SetGlobalLevel(logToZerologMap[level])
	}
}

func WithContextAlias(alias string) OptionSetter {
	return func(l *logger) {
		l.contextAlias = alias
	}
}

// Logger methods:

func (l *logger) SetMinLevel(level LogLevel) Logger {
	l.minLevel = level
	return l
}

func (l *logger) SetContext(ctx context.Context) Logger {
	// TODO: fix this
	// l.logger = l.logger.WithContext(ctx)
	return l
}

func (l *logger) SetKeyword(keyword string) Logger {
	l.contextAlias = keyword
	return l
}

func (l *logger) logEvent(level LogLevel, msg string, fields LogFields) Logger {
	if l.disabled || level < l.minLevel {
		return l
	}

	event := l.logger.WithLevel(logToZerologMap[level])
	if l.contextAlias != "" {
		fields["context"] = l.contextAlias
	}
	if fields != nil {
		event = event.Fields(fields)
	}
	if level >= WarnLevel {
		event = event.Stack()
	}
	event.Msg(msg)
	return l
}

func (l *logger) Trace(msg string, fields LogFields) Logger {
	return l.logEvent(TraceLevel, msg, fields)
}

func (l *logger) Debug(msg string, fields LogFields) Logger {
	return l.logEvent(DebugLevel, msg, fields)
}

func (l *logger) Info(msg string, fields LogFields) Logger {
	return l.logEvent(InfoLevel, msg, fields)
}

func (l *logger) Warn(msg string, fields LogFields) Logger {
	return l.logEvent(WarnLevel, msg, fields)
}

func (l *logger) Error(msg string, fields LogFields) Logger {
	return l.logEvent(ErrorLevel, msg, fields)
}

func (l *logger) Fatal(msg string, fields LogFields) Logger {
	return l.logEvent(FatalLevel, msg, fields)
}

func (l *logger) Panic(msg string, fields LogFields) Logger {
	return l.logEvent(PanicLevel, msg, fields)
}

func (l *logger) Disable() Logger {
	l.disabled = true
	return l
}

func (l *logger) Println(msg string) Logger {
	if !l.disabled {
		l.logger.Print(msg)
	}
	return l
}
