package glog

import (
	"context"
	"fmt"
)

// ContextToLogger returns the attached logger if available
func ContextToLogger(ctx context.Context) (Logger, error) {
	l, ok := ctx.Value(logStr("logger")).(Logger)
	if ok {
		return l, nil
	}

	return nil, fmt.Errorf("logger not found")
}
