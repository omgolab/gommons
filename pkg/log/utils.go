package glog

import (
	"context"
	"fmt"
)

// LoggerToContext returns a context with the attached logger
func LoggerToContext[T Logger](parentCtx context.Context, l T) context.Context {
	return context.WithValue(parentCtx, LogStr("logger"), l)
}

// ContextToLogger returns the attached logger if available
func ContextToLogger[T any](ctx context.Context) (T, error) {
	l, ok := ctx.Value(LogStr("logger")).(T)
	if ok {
		return l, nil
	}

	return l, fmt.Errorf("logger not found")
}
