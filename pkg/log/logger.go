package glog

import (
	"fmt"
	"io"
	"sync"

	gcollections "github.com/omgolab/go-commons/pkg/collections"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type LogFields map[string]any
type LogLevel int
type LogStr string

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

var (
	logToZerologMap = map[LogLevel]zerolog.Level{
		NoLevel:    zerolog.NoLevel,
		TraceLevel: zerolog.TraceLevel,
		DebugLevel: zerolog.DebugLevel,
		InfoLevel:  zerolog.InfoLevel,
		WarnLevel:  zerolog.WarnLevel,
		ErrorLevel: zerolog.ErrorLevel,
		FatalLevel: zerolog.FatalLevel,
		PanicLevel: zerolog.PanicLevel,
	}

	globalInitOnce sync.Once
)

type Logger interface {
	Event(msg string, level LogLevel, err error, csfCount int, fields ...LogFields)
	Trace(msg string, fields ...LogFields)
	Debug(msg string, fields ...LogFields)
	Info(msg string, fields ...LogFields)
	Warn(msg string, fields ...LogFields)
	Error(msg string, err error, fields ...LogFields)
	Fatal(msg string, err error, fields ...LogFields)
	Panic(msg string, err error, fields ...LogFields)
	Println(msg ...any)
	Printf(format string, v ...any)
	SetMinGlobalLogLevel(minLevel LogLevel) Logger
	SetMinCallerAttachLevel(minLevel LogLevel) Logger
	SetContextNS(keyword string) Logger
	DisableStackTraceOnError() Logger
	DisableTimestamp() Logger
	DisableAllLoggers() Logger
	update(nuc uniqueCfg) Logger
}

type sharedCfg struct {
	mu          sync.RWMutex
	minLogLevel LogLevel
	isDisabled  bool
	writers     []io.Writer
	timeFormat  string
}

type uniqueCfg struct {
	isTimestampOff  bool
	isStackTraceOff bool
	minCallerLevel  LogLevel
	ns              string
}

type logCfg struct {
	sc *sharedCfg

	mu sync.RWMutex
	zl zerolog.Logger
	uc uniqueCfg
}

func New(options ...LogOption) (Logger, error) {
	globalInitOnce.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	})

	defaultWriter := newConsoleWriter("Mon 02-Jan-06 03:04:05 PM -0700")
	l := &logCfg{
		sc: &sharedCfg{
			minLogLevel: DebugLevel,
			writers:     []io.Writer{defaultWriter},
			timeFormat:  defaultWriter.TimeFormat,
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

	if err := l.rebuildLogger(); err != nil {
		return nil, err
	}

	return l.update(l.uc), nil
}

func newConsoleWriter(format string) *zerolog.ConsoleWriter {
	cw := zerolog.NewConsoleWriter()
	cw.TimeFormat = format
	return &cw
}

func (l *logCfg) rebuildLogger() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rebuildLoggerLocked()
}

func (l *logCfg) rebuildLoggerLocked() error {
	writers := l.writersSnapshot()
	if len(writers) == 0 {
		return fmt.Errorf("glog: logger has no writers configured")
	}

	base := zerolog.New(zerolog.MultiLevelWriter(writers...))
	ctx := base.With()
	if l.uc.ns != "" {
		ctx = ctx.Str("context-ns", l.uc.ns)
	}
	if !l.uc.isTimestampOff {
		ctx = ctx.Timestamp()
	}
	if !l.uc.isStackTraceOff {
		ctx = ctx.Stack()
	}

	l.zl = ctx.Logger()
	return nil
}

func (l *logCfg) writersSnapshot() []io.Writer {
	l.sc.mu.RLock()
	defer l.sc.mu.RUnlock()
	writers := make([]io.Writer, len(l.sc.writers))
	copy(writers, l.sc.writers)
	return writers
}

func (l *logCfg) snapshot() (zerolog.Logger, uniqueCfg) {
	l.mu.RLock()
	logger := l.zl
	uc := l.uc
	l.mu.RUnlock()
	return logger, uc
}

func (l *logCfg) shouldSkip(level LogLevel) bool {
	l.sc.mu.RLock()
	disabled := l.sc.isDisabled || level < l.sc.minLogLevel
	l.sc.mu.RUnlock()
	return disabled
}

func (l *logCfg) update(nuc uniqueCfg) Logger {
	l.mu.Lock()
	l.uc = nuc
	l.rebuildLoggerLocked()
	l.mu.Unlock()
	return l
}

func (l *logCfg) SetMinGlobalLogLevel(minLevel LogLevel) Logger {
	l.sc.mu.Lock()
	l.sc.minLogLevel = minLevel
	l.sc.mu.Unlock()
	return l
}

func (l *logCfg) DisableStackTraceOnError() Logger {
	nuc := l.uc
	nuc.isStackTraceOff = true
	return l.update(nuc)
}

func (l *logCfg) DisableTimestamp() Logger {
	nuc := l.uc
	nuc.isTimestampOff = true
	return l.update(nuc)
}

func (l *logCfg) SetMinCallerAttachLevel(minLevel LogLevel) Logger {
	nuc := l.uc
	nuc.minCallerLevel = minLevel
	return l.update(nuc)
}

func (l *logCfg) SetContextNS(keyword string) Logger {
	nuc := l.uc
	nuc.ns = keyword
	return l.update(nuc)
}

func (l *logCfg) DisableAllLoggers() Logger {
	l.sc.mu.Lock()
	l.sc.isDisabled = true
	l.sc.mu.Unlock()
	return l
}

func (l *logCfg) Event(msg string, level LogLevel, err error, csfCount int, fields ...LogFields) {
	if l.shouldSkip(level) {
		return
	}

	logger, uc := l.snapshot()
	event := logger.WithLevel(logToZerologMap[level])

	if len(fields) > 0 {
		if merged := gcollections.MergeMaps(fields...); merged != nil {
			event = event.Fields(merged)
		}
	}

	if err != nil {
		event = event.Err(err)
	}

	if uc.minCallerLevel <= level {
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
	if l.shouldSkip(DebugLevel) {
		return
	}
	l.Event(fmt.Sprint(msg...), DebugLevel, nil, 3)
}

func (l *logCfg) Printf(format string, v ...any) {
	if l.shouldSkip(DebugLevel) {
		return
	}
	l.Event(fmt.Sprintf(format, v...), DebugLevel, nil, 3)
}

func (l *logCfg) SetOutput(w io.Writer) {
	l.sc.mu.Lock()
	l.sc.writers = []io.Writer{w}
	l.sc.mu.Unlock()
	_ = l.rebuildLogger()
}

func (l *logCfg) SetPrefix(string) {}

func (l *logCfg) Flags() int { return 0 }

func (l *logCfg) SetFlags(int) {}

func (l *logCfg) Writer() io.Writer {
	l.sc.mu.RLock()
	writers := make([]io.Writer, len(l.sc.writers))
	copy(writers, l.sc.writers)
	l.sc.mu.RUnlock()
	return zerolog.MultiLevelWriter(writers...)
}
