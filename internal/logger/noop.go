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
func (l *NoopLogger) Debug(_ string, _ ...interface{}) {}

// Info implements Logger
func (l *NoopLogger) Info(_ string, _ ...interface{}) {}

// Warn implements Logger
func (l *NoopLogger) Warn(_ string, _ ...interface{}) {}

// Error implements Logger
func (l *NoopLogger) Error(_ string, _ ...interface{}) {}

// Fatal implements Logger
func (l *NoopLogger) Fatal(_ string, _ ...interface{}) {
	os.Exit(1)
}

// WithError implements Logger
func (l *NoopLogger) WithError(_ error) Logger {
	return l
}

// WithField implements Logger
func (l *NoopLogger) WithField(_ string, _ interface{}) Logger {
	return l
}

// WithFields implements Logger
func (l *NoopLogger) WithFields(_ map[string]interface{}) Logger {
	return l
}
