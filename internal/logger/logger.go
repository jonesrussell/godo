package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonesrussell/godo/internal/config"
	"go.uber.org/zap"
)

var log *zap.Logger

// InitializeWithConfig sets up the logger with the provided configuration
func InitializeWithConfig(cfg config.LoggingConfig) error {
	zapConfig := zap.NewProductionConfig()

	// Configure log level
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	zapConfig.Level = level

	// Ensure log directories exist
	if err := ensureLogDirectories(cfg.Output); err != nil {
		return fmt.Errorf("failed to create log directories: %w", err)
	}

	// Configure output paths
	zapConfig.OutputPaths = cfg.Output
	zapConfig.ErrorOutputPaths = cfg.ErrorOutput

	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("failed to build logger: %w", err)
	}

	log = logger
	return nil
}

// ensureLogDirectories creates directories for log files if they don't exist
func ensureLogDirectories(paths []string) error {
	for _, path := range paths {
		if path == "stdout" || path == "stderr" {
			continue
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	if log != nil {
		log.Sugar().Debugf(msg, args...)
	}
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	if log != nil {
		log.Sugar().Infof(msg, args...)
	}
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	if log != nil {
		log.Sugar().Warnf(msg, args...)
	}
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	if log != nil {
		log.Sugar().Errorf(msg, args...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	if log != nil {
		log.Sugar().Fatalf(msg, args...)
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
