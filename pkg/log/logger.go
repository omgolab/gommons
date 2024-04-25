package gclog

import (
	"fmt"
	"io"

	gccollections "github.com/omgolab/go-commons/pkg/collections"
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
	// log methods
	Event(msg string, level LogLevel, err error, csfCount int, fields ...LogFields)
	Trace(msg string, fields ...LogFields)
	Debug(msg string, fields ...LogFields)
	Info(msg string, fields ...LogFields)
	Warn(msg string, fields ...LogFields)
	Error(msg string, err error, fields ...LogFields)
	Fatal(msg string, err error, fields ...LogFields)
	Panic(msg string, err error, fields ...LogFields)
	Println(msg ...any)
	Printf(format string, v ...interface{})
	// settings methods
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
}

// uniqueCfg separate for each logger instance
type uniqueCfg struct {
	isTimestampOff  bool     // default ON
	isStackTraceOff bool     // default ON
	minCallerLevel  LogLevel // default WarnLevel
	context         string   // default context (namespace)
}

type logCfg struct {
	sc *sharedCfg
	zl zerolog.Logger
	uc uniqueCfg
}

func New(options ...LogOption) (Logger, error) {
	l := &logCfg{
		sc: &sharedCfg{
			minLogLevel: DebugLevel,
			writers:     []io.Writer{zerolog.NewConsoleWriter()},
			timeFormat:  "Mon 02-Jan-06 03:04:05 PM -0700",
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

	// update the zero logger
	l.zl = zerolog.New(zerolog.MultiLevelWriter(l.sc.writers...))
	zerolog.TimeFieldFormat = l.sc.timeFormat

	return l.update(l.uc), nil
}

// Logger methods:

func (l *logCfg) update(nuc uniqueCfg) Logger {
	// create a copy of the logger if the same config is not provided
	if nuc != l.uc {
		ll := *l
		l = &ll
	}

	ctx := l.zl.With()

	// 1. update the context
	if nuc.context != "" {
		if l.uc.context != "" {
			// we need to clear the old context
			ctx = zerolog.New(zerolog.MultiLevelWriter(l.sc.writers...)).With()
		}
		ctx = ctx.Str("context", nuc.context)
	}

	// 2. update timestamp
	if !nuc.isTimestampOff {
		ctx = ctx.Timestamp()
	}

	// 3. update stack trace on error
	if !nuc.isStackTraceOff {
		ctx = ctx.Stack()
	}

	// update the logger
	l.zl = ctx.Logger()
	l.uc = nuc

	return l
}

func (l *logCfg) SetMinGlobalLogLevel(minLevel LogLevel) Logger {
	l.sc.minLogLevel = minLevel
	return l
}

func (l *logCfg) DisableStackTraceOnError() Logger {
	// copy current unique config
	nuc := l.uc
	nuc.isStackTraceOff = true
	return l.update(nuc)
}

func (l *logCfg) DisableTimestamp() Logger {
	// copy current unique config
	nuc := l.uc
	nuc.isTimestampOff = true
	return l.update(nuc)
}

func (l *logCfg) SetMinCallerAttachLevel(minLevel LogLevel) Logger {
	// copy current unique config
	nuc := l.uc
	nuc.minCallerLevel = minLevel
	return l.update(nuc)
}

// SetContext returns a new child logger with the context set
func (l *logCfg) SetContext(keyword string) Logger {
	// copy current unique config
	nuc := l.uc
	nuc.context = keyword
	return l.update(nuc)
}

func (l *logCfg) DisableAllLoggers() Logger {
	l.sc.isDisabled = true
	return l
}

func (l *logCfg) Event(msg string, level LogLevel, err error, csfCount int, fields ...LogFields) {
	if l.sc.isDisabled || level < l.sc.minLogLevel {
		return
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
	if l.uc.minCallerLevel <= level {
		// csfCount = callerSkipFrameCount, normally to skip the caller
		event = event.Caller(csfCount)
	}

	event.Msg(msg)
}

func (l *logCfg) Trace(msg string, fields ...LogFields) {
	l.Event(msg, TraceLevel, nil, 2, fields...)
}

func (l *logCfg) Debug(msg string, fields ...LogFields) {
	l.Event(msg, DebugLevel, nil, 2, fields...)
}

func (l *logCfg) Info(msg string, fields ...LogFields) {
	l.Event(msg, InfoLevel, nil, 2, fields...)
}

func (l *logCfg) Warn(msg string, fields ...LogFields) {
	l.Event(msg, WarnLevel, nil, 2, fields...)
}

func (l *logCfg) Error(msg string, err error, fields ...LogFields) {
	l.Event(msg, ErrorLevel, err, 2, fields...)
}

func (l *logCfg) Fatal(msg string, err error, fields ...LogFields) {
	l.Event(msg, FatalLevel, err, 2, fields...)
}

func (l *logCfg) Panic(msg string, err error, fields ...LogFields) {
	l.Event(msg, PanicLevel, err, 2, fields...)
}

func (l *logCfg) Println(msg ...any) {
	if !l.sc.isDisabled {
		l.zl.Log().CallerSkipFrame(1).Msg(fmt.Sprint(msg...))
	}
}

func (l *logCfg) Printf(format string, v ...interface{}) {
	if !l.sc.isDisabled {
		l.zl.Log().CallerSkipFrame(1).Msg(fmt.Sprintf(format, v...))
	}
}
