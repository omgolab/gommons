package gccustlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	gclog "github.com/omar391/go-commons/pkg/log"
	"github.com/rs/zerolog"
)

type TaggedLogger interface {
	gclog.Logger
	LogTag(msg string, level gclog.LogLevel, err error, csfCount int, fields ...gclog.LogFields)
	GetDelimiter() string
	IsTimestampFormatterEnabled() bool
	IsLevelFormatterEnabled() bool
	IsCallerFormatterEnabled() bool
	IsErrorFormatterEnabled() bool
	// UpdateBaseLogger(l gclog.Logger)
	// add others if needed
}

type tlCfg struct {
	gclog.Logger
	tag                    string
	delimiter              string
	writer                 io.Writer
	partsOrder             []string
	timestampFormatter     zerolog.Formatter
	levelFormatter         zerolog.Formatter
	messageFormatter       zerolog.Formatter
	callerFormatter        zerolog.Formatter
	fieldNameFormatter     zerolog.Formatter
	fieldValueFormatter    zerolog.Formatter
	errFieldNameFormatter  zerolog.Formatter
	errFieldValueFormatter zerolog.Formatter
	logOpts                []gclog.LogOption
}

func (fw *tlCfg) GetDelimiter() string {
	return fw.delimiter
}

// func (fw *tlCfg) UpdateBaseLogger(l gclog.Logger) {
// 	fw.Logger = l
// }

func (fw *tlCfg) IsTimestampFormatterEnabled() bool {
	return fw.timestampFormatter("") != "-x-"
}

func (fw *tlCfg) IsLevelFormatterEnabled() bool {
	return fw.levelFormatter("") != "-x-"
}

func (fw *tlCfg) IsCallerFormatterEnabled() bool {
	return fw.callerFormatter("") != "-x-"
}

func (fw *tlCfg) IsErrorFormatterEnabled() bool {
	return fw.errFieldNameFormatter("") != "-x-"
}

func (fw *tlCfg) LogTag(msg string, level gclog.LogLevel, err error, csfCount int, fields ...gclog.LogFields) {
	msg += fw.delimiter + " " + fw.tag
	fw.Event(msg, level, err, csfCount, fields...)
}

// Write implements io.Writer and only writes to underlying writer if the tag exists
func (fw *tlCfg) Write(b []byte) (n int, err error) {
	ln := len(b) // used for returning the original length
	bi := []byte(fw.tag)
	l := len(bi) + 1 // 1 for the delimiter or a ending newline

	for {
		m := bytes.LastIndex(b, bi)
		if m < 0 {
			break
		}

		j := m + l
		if j < len(b) && b[j] == ' ' {
			j++
		} else if j == len(b) {
			j--
		}

		copy(b[m:], b[j:])
		b = b[:len(b)-(j-m)]
	}

	bl := len(b)
	if bl < ln {
		d := 0
		if bl >= 2 {
			if b[bl-3] == fw.delimiter[0] {
				d = 3
			} else if b[bl-2] == fw.delimiter[0] {
				d = 2
			}
		}
		// remove ending delimiter
		if d > 0 {
			b[bl-d] = '\n'
			b = b[:bl-d+1]
		}

		_, err = fw.writer.Write(b)
	}

	return ln, err
}

func New(tag string, w io.Writer, tlOpts ...TaggedLoggerOption) (l TaggedLogger, err error) {
	if tag == "" || w == nil {
		return nil, errors.New("tag or writer is invalid")
	}
	defaultFormatter := func(i interface{}) string { return "-x-" }

	cfg := &tlCfg{
		tag:       fmt.Sprintf("tag=%s-%d", tag, time.Now().UnixNano()), // make the tag unique
		delimiter: ",",
		writer:    w,
		partsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
		fieldNameFormatter:  defaultFormatter,
		fieldValueFormatter: defaultFormatter,
	}

	// enable default formatters
	dfo := []TaggedLoggerOption{
		WithTimestampFormatter(),
		WithMessageFormatter(),
		WithErrFieldNameFormatter(),
		WithErrFieldValueFormatter(),
		WithCallerFormatter(),
		WithLevelFormatter(),
	}
	tlOpts = append(dfo, tlOpts...)

	// apply the options
	for _, opt := range tlOpts {
		opt(cfg)
	}

	// set console writer
	cw := zerolog.ConsoleWriter{
		Out:                 cfg,
		FormatTimestamp:     cfg.timestampFormatter,
		FormatLevel:         cfg.levelFormatter,
		FormatCaller:        cfg.callerFormatter,
		FormatMessage:       cfg.messageFormatter,
		FormatFieldName:     cfg.fieldNameFormatter,
		FormatFieldValue:    cfg.fieldValueFormatter,
		FormatErrFieldName:  cfg.errFieldNameFormatter,
		FormatErrFieldValue: cfg.errFieldValueFormatter,
		NoColor:             true,
		PartsOrder:          cfg.partsOrder,
	}
	logOpts := append(cfg.logOpts, gclog.WithPrimaryWriter(cw))

	cfg.Logger, err = gclog.New(logOpts...)
	return cfg, err
}
