package gcstrlog_test

import (
	"errors"
	"testing"

	"github.com/tj/assert"

	sl "github.com/omar391/go-commons/pkg/log/gccustlog/string"
)

func TestNew(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		arrayLogger, err := sl.New()

		assert.NoError(t, err)
		assert.NotNil(t, arrayLogger)
	})
}

func TestStringLogger(t *testing.T) {

	t.Run("append and retrieve string logs", func(t *testing.T) {
		strLogger, _ := sl.New()
		strLogger.AppendStringErr("test message", errors.New("Test err"))

		logs := strLogger.GetStringLogs()
		assert.Equal(t, 1, len(logs))
		assert.Contains(t, logs[0], "test message")
	})

	t.Run("append multiple string logs", func(t *testing.T) {
		strLogger, _ := sl.New()
		strLogger.AppendString("first message")
		strLogger.AppendString("second message")

		logs := strLogger.GetStringLogs()
		assert.Equal(t, 2, len(logs))
		assert.Contains(t, logs[0], "first message")
		assert.Contains(t, logs[1], "second message")
	})
}
