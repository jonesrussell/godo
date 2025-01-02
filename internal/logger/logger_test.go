package logger

import "testing"

func TestLoggerBasicOperations(t *testing.T) {
	t.Run("Debug logs message", func(t *testing.T) {
		logger := NewTestLogger(t)
		logger.Debug("test debug message")
		// Add assertion if needed
	})

	t.Run("Info logs message", func(t *testing.T) {
		logger := NewTestLogger(t)
		logger.Info("test info message")
		// Add assertion if needed
	})

	t.Run("Warn logs message", func(t *testing.T) {
		logger := NewTestLogger(t)
		logger.Warn("test warn message")
		// Add assertion if needed
	})

	t.Run("Error logs message", func(t *testing.T) {
		logger := NewTestLogger(t)
		logger.Error("test error message")
		// Add assertion if needed
	})
}
