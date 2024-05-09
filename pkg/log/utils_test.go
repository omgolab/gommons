package glog_test

import (
	"context"
	"testing"

	glog "github.com/omgolab/go-commons/pkg/log"
	gcustomlog "github.com/omgolab/go-commons/pkg/log/custom"
	gstrlog "github.com/omgolab/go-commons/pkg/log/custom/string"
)

// Returns the attached logger if available with correct type.
func TestContextToLogger_ReturnsLoggerIfAvailable(t *testing.T) {
	l, _ := gstrlog.New([]gcustomlog.FilterOption{})
	ctx := context.WithValue(context.Background(), glog.LogStr("logger"), l)
	logger, err := glog.ContextToLogger[gstrlog.StringLogger](ctx)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if logger != l {
		t.Errorf("Expected logger to be StringLogger, but got: %v", logger)
	}
}
