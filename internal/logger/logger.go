package logger

import (
	"fmt"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// InitializeWithConfig sets up the logger with the provided configuration
func InitializeWithConfig(cfg common.LogConfig) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	// Set log level from config
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// Initialize logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	log = logger
	return logger, nil
}

// Debug logs a debug message
func Debug(msg string, keysAndValues ...interface{}) {
	if log != nil {
		sugar := log.Sugar()
		sugar.Debugw(msg, keysAndValues...)
	}
}

// Info logs an info message
func Info(msg string, keysAndValues ...interface{}) {
	if log != nil {
		sugar := log.Sugar()
		sugar.Infow(msg, keysAndValues...)
	}
}

// Error logs an error message
func Error(msg string, keysAndValues ...interface{}) {
	if log != nil {
		sugar := log.Sugar()
		sugar.Errorw(msg, keysAndValues...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, keysAndValues ...interface{}) {
	if log != nil {
		sugar := log.Sugar()
		sugar.Fatalw(msg, keysAndValues...)
	}
}

// Format returns a formatted string using the provided format and args
func Format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
