package gcstrlog

import (
	gccustomlog "github.com/omar391/go-commons/pkg/log/gccustlog"
)

type StringLoggerOption func(*arrayCfg)

func WithInitialDataSize(size int) StringLoggerOption {
	return func(cfg *arrayCfg) {
		cfg.initialDataSize = size
	}
}

func WithTaggedLoggerOptions(opts ...gccustomlog.TaggedLoggerOption) StringLoggerOption {
	return func(cfg *arrayCfg) {
		cfg.tlOpts = opts
	}
}
