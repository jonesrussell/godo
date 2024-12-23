package logger

import "go.uber.org/zap"

// TestLogger returns a logger suitable for testing
func NewTestLogger() Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, _ := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return &ZapLogger{
		log: logger.Sugar(),
	}
}

// NoopLogger returns a logger that does nothing (for testing)
func NewNoopLogger() Logger {
	return &noopLogger{}
}

// noopLogger implements Logger but does nothing
type noopLogger struct{}

// Basic logging methods
func (l *noopLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l *noopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l *noopLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l *noopLogger) Error(msg string, keysAndValues ...interface{}) {}
func (l *noopLogger) Fatal(msg string, keysAndValues ...interface{}) {}

// Helper methods
func (l *noopLogger) WithField(key string, value interface{}) Logger {
	return l
}

func (l *noopLogger) WithFields(fields map[string]interface{}) Logger {
	return l
}

func (l *noopLogger) WithError(err error) Logger {
	return l
}
