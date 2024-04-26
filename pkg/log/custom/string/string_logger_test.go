package gstrlog_test

import (
	"errors"
	"testing"

	log "github.com/omgolab/go-commons/pkg/log"
	"github.com/tj/assert"

	filter "github.com/omgolab/go-commons/pkg/log/custom"
	gstrlog "github.com/omgolab/go-commons/pkg/log/custom/string"
)

func TestNew(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		filterOpts := []filter.FilterOption{}
		logOpts := []log.LogOption{}

		arrayLogger, err := gstrlog.New(filterOpts, logOpts...)

		assert.NoError(t, err)
		assert.NotNil(t, arrayLogger)
	})
}

func TestStringLogger(t *testing.T) {
	filterOpts := []filter.FilterOption{}
	logOpts := []log.LogOption{}
	strLogger, _ := gstrlog.New(filterOpts, logOpts...)

	t.Run("append and retrieve string logs", func(t *testing.T) {
		strLogger.AppendStringErr("test message", errors.New("Test err"))

		logs := strLogger.GetStringLogs()
		assert.Equal(t, 1, len(logs))
		assert.Contains(t, logs[0], "test message")
	})

	t.Run("append multiple string logs", func(t *testing.T) {
		strLogger.AppendString("first message")
		strLogger.AppendString("second message")

		logs := strLogger.GetStringLogs()
		assert.Equal(t, 3, len(logs)) // Note that there is already 1 log from the previous sub-test
		assert.Contains(t, logs[1], "first message")
		assert.Contains(t, logs[2], "second message")
	})
}
