package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerOperations(t *testing.T) {
	cfg := &Config{
		Level:       "debug",
		Development: true,
		Encoding:    "console",
	}
	logger, _, err := NewLogger(cfg)
	if err != nil {
		t.Errorf("Failed to create test logger: %v", err)
		t.FailNow()
	}

	t.Run("Debug logs message", func(t *testing.T) {
		logger.Debug("test debug message", "key", "value")
	})

	t.Run("Info logs message", func(t *testing.T) {
		logger.Info("test info message", "key", "value")
	})

	t.Run("Warn logs message", func(t *testing.T) {
		logger.Warn("test warn message", "key", "value")
	})

	t.Run("Error logs message", func(t *testing.T) {
		logger.Error("test error message", "key", "value")
	})

	t.Run("WithError adds error field", func(t *testing.T) {
		err := errors.New("test error")
		errorLogger := logger.WithError(err)
		assert.NotNil(t, errorLogger)
		errorLogger.Error("error occurred")
	})
}
