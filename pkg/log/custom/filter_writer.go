package gcustomlog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	glog "github.com/omgolab/go-commons/pkg/log"
	"github.com/rs/zerolog"
)

type FilterLogger interface {
	glog.Logger
	TagLog(msg string, level glog.LogLevel, err error, csfCount int, fields ...glog.LogFields)
	GetDelimiter() string
	UpdateBaseLogger(l glog.Logger)
	IsTimestampFormatterEnabled() bool
	IsLevelFormatterEnabled() bool
	IsCallerFormatterEnabled() bool
	IsErrorFormatterEnabled() bool
	// add others if needed
}

type filterWriter struct {
	glog.Logger
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
}

func (fw *filterWriter) GetDelimiter() string {
	return fw.delimiter
}

func (fw *filterWriter) UpdateBaseLogger(l glog.Logger) {
	fw.Logger = l
}

func (fw *filterWriter) IsTimestampFormatterEnabled() bool {
	return fw.timestampFormatter("") != ""
}

func (fw *filterWriter) IsLevelFormatterEnabled() bool {
	return fw.levelFormatter("") != ""
}

func (fw *filterWriter) IsCallerFormatterEnabled() bool {
	return fw.callerFormatter("") != ""
}

func (fw *filterWriter) IsErrorFormatterEnabled() bool {
	return fw.errFieldNameFormatter("") != ""
}

func (fw *filterWriter) TagLog(msg string, level glog.LogLevel, err error, csfCount int, fields ...glog.LogFields) {
	msg += fw.delimiter + " " + fw.tag
	fw.Event(msg, level, err, csfCount, fields...)
}

// Write implements io.Writer and only writes to underlying writer if the tag exists
func (fw *filterWriter) Write(b []byte) (n int, err error) {
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

func New(tag string, w io.Writer, filterOpts []FilterOption, logOpts ...glog.LogOption) (l FilterLogger, err error) {
	if tag == "" || w == nil {
		return nil, errors.New("tag or writer is invalid")
	}
	defaultFormatter := func(i interface{}) string { return "" }

	fw := &filterWriter{
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
	dfo := []FilterOption{
		WithTimestampFormatter(),
		WithMessageFormatter(),
		WithErrFieldNameFormatter(),
		WithErrFieldValueFormatter(),
		WithCallerFormatter(),
		WithLevelFormatter(),
	}
	filterOpts = append(dfo, filterOpts...)

	// apply the options
	for _, opt := range filterOpts {
		if err = opt(fw); err != nil {
			return nil, err
		}
	}

	// set console writer
	cw := zerolog.ConsoleWriter{
		Out:                 fw,
		FormatTimestamp:     fw.timestampFormatter,
		FormatLevel:         fw.levelFormatter,
		FormatCaller:        fw.callerFormatter,
		FormatMessage:       fw.messageFormatter,
		FormatFieldName:     fw.fieldNameFormatter,
		FormatFieldValue:    fw.fieldValueFormatter,
		FormatErrFieldName:  fw.errFieldNameFormatter,
		FormatErrFieldValue: fw.errFieldValueFormatter,
		NoColor:             true,
		PartsOrder:          fw.partsOrder,
	}
	logOpts = append(logOpts, glog.WithMultiLogger(cw))

	fw.Logger, err = glog.New(logOpts...)
	return fw, err
}
