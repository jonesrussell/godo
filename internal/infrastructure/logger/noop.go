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
func (l *NoopLogger) Debug(_ string, _ ...any) {}

// Info implements Logger
func (l *NoopLogger) Info(_ string, _ ...any) {}

// Warn implements Logger
func (l *NoopLogger) Warn(_ string, _ ...any) {}

// Error implements Logger
func (l *NoopLogger) Error(_ string, _ ...any) {}

// Fatal implements Logger
func (l *NoopLogger) Fatal(_ string, _ ...any) {
	os.Exit(1)
}

// WithError implements Logger
func (l *NoopLogger) WithError(_ error) Logger {
	return l
}

// WithField implements Logger
func (l *NoopLogger) WithField(_ string, _ any) Logger {
	return l
}

// WithFields implements Logger
func (l *NoopLogger) WithFields(_ map[string]any) Logger {
	return l
}
