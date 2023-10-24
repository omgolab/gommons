package gccustomlog

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type FilterOption func(*filterWriter) error

func WithDelimiter(delim rune) FilterOption {
	return func(fw *filterWriter) error {
		fw.delimiter = string(delim)
		return nil
	}
}

func WithTimestampFormatter(f ...string) FilterOption {
	return func(fw *filterWriter) error {
		if len(f) == 0 || f[0] == "" {
			fw.timestampFormatter = func(i interface{}) string {
				return i.(string) + fw.delimiter
			}
			return nil
		}

		fw.timestampFormatter = func(i interface{}) string {
			v := i.(string)
			ts, err := time.ParseInLocation(f[0], v, time.Local)
			if err != nil {
				return v + fw.delimiter
			} else {
				return ts.Local().Format(time.DateTime) + fw.delimiter
			}
		}
		return nil
	}
}

func WithLevelFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		if len(f) == 0 || f[0] == nil {
			fw.levelFormatter = func(i interface{}) string {
				return strings.ToUpper(i.(string))[0:3] + fw.delimiter
			}
			return nil
		}

		fw.levelFormatter = f[0]
		return nil
	}
}

func WithMessageFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		if len(f) == 0 || f[0] == nil {
			fw.messageFormatter = func(i interface{}) string {
				s, b := i.(string)
				if i == nil || !b {
					return ""
				}
				return s + fw.delimiter
			}
			return nil
		}

		fw.messageFormatter = f[0]
		return nil
	}
}

func WithCallerFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		if len(f) == 0 || f[0] == nil {
			fw.callerFormatter = func(i interface{}) string {
				var s string
				var ok bool
				if s, ok = i.(string); !ok {
					return ""
				}
				if len(s) > 0 {
					if cwd, err := os.Getwd(); err == nil {
						if rel, err := filepath.Rel(cwd, s); err == nil {
							s = rel
						}
					}
				}
				return s + " >" + fw.delimiter
			}
			return nil
		}

		fw.callerFormatter = f[0]
		return nil
	}
}

func getDefaultFieldNameFn(f []zerolog.Formatter) zerolog.Formatter {
	if len(f) == 0 || f[0] == nil {
		return func(i interface{}) string {
			s, b := i.(string)
			if i == nil || !b {
				return ""
			}
			return s + "="
		}
	}

	return f[0]
}

func getDefaultFieldValueFn(f []zerolog.Formatter, d string) zerolog.Formatter {
	if len(f) == 0 || f[0] == nil {
		return func(i interface{}) string {
			s, b := i.(string)
			if i == nil || !b {
				return ""
			}
			return s + d
		}
	}

	return f[0]
}

func WithFieldNameFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		fw.fieldNameFormatter = getDefaultFieldNameFn(f)
		return nil
	}
}

func WithFieldValueFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		fw.fieldValueFormatter = getDefaultFieldValueFn(f, fw.delimiter)
		return nil
	}
}

func WithErrFieldNameFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		fw.errFieldNameFormatter = getDefaultFieldNameFn(f)
		return nil
	}
}

func WithErrFieldValueFormatter(f ...zerolog.Formatter) FilterOption {
	return func(fw *filterWriter) error {
		fw.errFieldValueFormatter = getDefaultFieldValueFn(f, fw.delimiter)
		return nil
	}
}

func WithPartsOrder(o []string) FilterOption {
	return func(fw *filterWriter) error {
		for _, v := range o {
			if !slices.Contains[[]string](fw.partsOrder, v) {
				fw.partsOrder = append(fw.partsOrder, o...)
			}
		}
		return nil
	}
}
