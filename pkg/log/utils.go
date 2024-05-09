package glog

import (
	"context"
	"fmt"
)

// ContextToLogger returns the attached logger if available
func ContextToLogger[T any](ctx context.Context) (T, error) {
	l, ok := ctx.Value(LogStr("logger")).(T)
	if ok {
		return l, nil
	}

	return l, fmt.Errorf("logger not found")
}
