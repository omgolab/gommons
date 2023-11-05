package gccustlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	gclog "github.com/omar391/go-commons/pkg/log"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

type TaggedLoggerOption func(*tlCfg)

func WithDelimiter(delim rune) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.delimiter = string(delim)
	}
}

func WithLoggerOptions(opts ...gclog.LogOption) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.logOpts = opts
	}
}

func WithTimestampFormatter(f ...string) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		if len(f) == 0 || f[0] == "" {
			cfg.timestampFormatter = func(i interface{}) string {
				return i.(string) + cfg.delimiter
			}
			return
		}

		cfg.timestampFormatter = func(i interface{}) string {
			v := i.(string)
			ts, err := time.ParseInLocation(f[0], v, time.Local)
			if err != nil {
				return v + cfg.delimiter
			} else {
				return ts.Local().Format(time.DateTime) + cfg.delimiter
			}
		}
	}
}

func WithLevelFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		if len(f) == 0 || f[0] == nil {
			cfg.levelFormatter = func(i interface{}) string {
				s, b := i.(string)
				if i == nil || !b || i == "" {
					return "???" + cfg.delimiter
				}
				switch s {
				case zerolog.LevelTraceValue:
					s = "TRC"
				case zerolog.LevelDebugValue:
					s = "DBG"
				case zerolog.LevelInfoValue:
					s = "INF"
				case zerolog.LevelWarnValue:
					s = "WRN"
				case zerolog.LevelErrorValue:
					s = "ERR"
				case zerolog.LevelFatalValue:
					s = "FTL"
				case zerolog.LevelPanicValue:
					s = "PNC"
				default:
					s = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
				}
				return s + cfg.delimiter
			}
			return
		}

		cfg.levelFormatter = f[0]
	}
}

func WithMessageFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		if len(f) == 0 || f[0] == nil {
			cfg.messageFormatter = func(i interface{}) string {
				s, b := i.(string)
				if i == nil || !b || i == "" {
					s = ""
				}
				return s + cfg.delimiter
			}
			return
		}

		cfg.messageFormatter = f[0]
	}
}

func WithCallerFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		if len(f) == 0 || f[0] == nil {
			cfg.callerFormatter = func(i interface{}) string {
				s, b := i.(string)
				if i == nil || !b || i == "" {
					return "" + cfg.delimiter
				}
				if len(s) > 0 {
					if cwd, err := os.Getwd(); err == nil {
						if rel, err := filepath.Rel(cwd, s); err == nil {
							s = rel
						}
					}
				}
				return s + " >" + cfg.delimiter
			}
			return
		}

		cfg.callerFormatter = f[0]
	}
}

func getDefaultFieldNameFn(f []zerolog.Formatter) zerolog.Formatter {
	if len(f) == 0 || f[0] == nil {
		return func(i interface{}) string {
			s, b := i.(string)
			if i == nil || !b || i == "" {
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
			if i == nil || !b || i == "" {
				return ""
			}
			return s + d
		}
	}

	return f[0]
}

func WithFieldNameFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.fieldNameFormatter = getDefaultFieldNameFn(f)
	}
}

func WithFieldValueFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.fieldValueFormatter = getDefaultFieldValueFn(f, cfg.delimiter)
	}
}

func WithErrFieldNameFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.errFieldNameFormatter = getDefaultFieldNameFn(f)
	}
}

func WithErrFieldValueFormatter(f ...zerolog.Formatter) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		cfg.errFieldValueFormatter = getDefaultFieldValueFn(f, cfg.delimiter)
	}
}

func WithPartsOrder(o []string) TaggedLoggerOption {
	return func(cfg *tlCfg) {
		for _, v := range o {
			if !slices.Contains[[]string](cfg.partsOrder, v) {
				cfg.partsOrder = append(cfg.partsOrder, o...)
			}
		}
	}
}
