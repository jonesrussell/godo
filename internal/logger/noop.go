package logger

import (
	"os"
)

// NoopLogger is a logger that does nothing
type NoopLogger struct{}

// NewNoopLogger creates a new logger that does nothing
func NewNoopLogger() Logger {
	return &NoopLogger{}
}

// Debug implements Logger
func (l *NoopLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Info implements Logger
func (l *NoopLogger) Info(msg string, keysAndValues ...interface{}) {}

// Warn implements Logger
func (l *NoopLogger) Warn(msg string, keysAndValues ...interface{}) {}

// Error implements Logger
func (l *NoopLogger) Error(msg string, keysAndValues ...interface{}) {}

// Fatal implements Logger
func (l *NoopLogger) Fatal(msg string, keysAndValues ...interface{}) {
	os.Exit(1)
}

// WithError implements Logger
func (l *NoopLogger) WithError(err error) Logger {
	return l
}

// WithField implements Logger
func (l *NoopLogger) WithField(key string, value interface{}) Logger {
	return l
}

// WithFields implements Logger
func (l *NoopLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}
