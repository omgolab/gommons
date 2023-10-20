package gclog

import (
	"fmt"
	"io"

	gccollections "github.com/omar391/go-commons/pkg/collections"
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
	Trace(msg string, fields ...LogFields) Logger
	Debug(msg string, fields ...LogFields) Logger
	Info(msg string, fields ...LogFields) Logger
	Warn(msg string, fields ...LogFields) Logger
	Error(msg string, err error, fields ...LogFields) Logger
	Fatal(msg string, err error, fields ...LogFields) Logger
	Panic(msg string, err error, fields ...LogFields) Logger
	Println(msg ...any) Logger
	Printf(format string, v ...interface{}) Logger
	SetMinGlobalLogLevel(minLevel LogLevel) Logger
	SetMinCallerAttachLevel(minLevel LogLevel) Logger
	SetContext(keyword string) Logger
	DisableStackTraceOnError() Logger
	DisableTimestamp() Logger
	DisableAllLoggers() Logger
	update(nuc uniqueCfg) Logger
}

// using sharedCfg so underlying data will be same on copy
type sharedCfg struct {
	minLogLevel LogLevel
	isDisabled  bool // default ON
	writers     []io.Writer
	timeFormat  string
	zlCtx       zerolog.Context
}

// uniqueCfg separate for each logger instance
type uniqueCfg struct {
	isTimestampOff  bool     // default ON
	isStackTraceOff bool     // default ON
	minCallerLevel  LogLevel // default WarnLevel
	context         string   // default context (namespace)
}

type logger struct {
	sc *sharedCfg
	zl zerolog.Logger
	uc uniqueCfg
}

func NewLogger(options ...OptionSetter) (Logger, error) {
	l := &logger{
		sc: &sharedCfg{
			minLogLevel: DebugLevel,
			writers:     []io.Writer{zerolog.NewConsoleWriter()},
			timeFormat:  "Mon, 02 Jan 06 5:04:05PM -0700",
		},
		uc: uniqueCfg{
			minCallerLevel: WarnLevel,
		},
	}

	for _, opt := range options {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	// 1. update the writers and base logger
	l.sc.zlCtx = zerolog.New(zerolog.MultiLevelWriter(l.sc.writers...)).With()
	// 2. update timestamp format
	zerolog.TimeFieldFormat = l.sc.timeFormat

	return l.update(l.uc), nil
}

// Logger methods:

func (l *logger) update(nuc uniqueCfg) Logger {
	// create a copy of the logger if the same config is not provided
	if nuc != l.uc {
		ll := *l
		l = &ll
	}

	ctx := l.sc.zlCtx

	// 1. update timestamp
	if !nuc.isTimestampOff {
		ctx = ctx.Timestamp()
	}

	// 2. update stack trace on error
	if !nuc.isStackTraceOff {
		ctx = ctx.Stack()
	}

	// 3. update the context
	if nuc.context != "" {
		ctx = ctx.Str("context", nuc.context)
	}

	// update the logger
	l.zl = ctx.Logger()

	return l
}

func (l *logger) SetMinGlobalLogLevel(minLevel LogLevel) Logger {
	l.sc.minLogLevel = minLevel
	return l
}

func (l *logger) DisableStackTraceOnError() Logger {
	// copy current unique config
	nuc := l.uc
	nuc.isStackTraceOff = true
	return l.update(nuc)
}

func (l *logger) DisableTimestamp() Logger {
	// copy current unique config
	nuc := l.uc
	nuc.isTimestampOff = true
	return l.update(nuc)
}

func (l *logger) SetMinCallerAttachLevel(minLevel LogLevel) Logger {
	// copy current unique config
	nuc := l.uc
	nuc.minCallerLevel = minLevel
	return l.update(nuc)
}

// SetContext returns a new child logger with the context set
func (l *logger) SetContext(keyword string) Logger {
	// copy current unique config
	nuc := l.uc
	nuc.context = keyword
	return l.update(nuc)
}

func (l *logger) logEvent(level LogLevel, msg string, err error, fields []LogFields) Logger {
	if l.sc.isDisabled || level < l.sc.minLogLevel {
		return l
	}

	event := l.zl.WithLevel(logToZerologMap[level])

	if len(fields) > 0 {
		lfs := gccollections.MergeMaps[LogFields](fields...)
		if lfs != nil {
			event = event.Fields(lfs)
		}
	}

	if err != nil {
		event = event.Err(err)
	}

	// use the caller level if enabled
	if l.uc.minCallerLevel >= level {
		event = event.Caller()
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

func (l *logger) DisableAllLoggers() Logger {
	l.sc.isDisabled = true
	return l
}

func (l *logger) Println(msg ...any) Logger {
	if !l.sc.isDisabled {
		l.zl.Log().CallerSkipFrame(1).Msg(fmt.Sprint(msg...))
	}
	return l
}

func (l *logger) Printf(format string, v ...interface{}) Logger {
	if !l.sc.isDisabled {
		l.zl.Log().CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
	}
	return l
}
