package gccsvlog

import "github.com/omar391/go-commons/pkg/log/gccustlog"

type CsvOption func(*csvCfg)

func WithTruncateOnHeadersMissing() CsvOption {
	return func(ch *csvCfg) {
		ch.truncateOnHeadersMissing = true
	}
}

func WithTaggedLoggerOptions(opts ...gccustlog.TaggedLoggerOption) CsvOption {
	return func(cfg *csvCfg) {
		cfg.tlOpts = opts
	}
}

func WithHeaders(headers []string) CsvOption {
	return func(cfg *csvCfg) {
		cfg.headers = headers
	}
}
