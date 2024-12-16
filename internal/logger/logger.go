package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

var log *zap.Logger

// Initialize sets up the logger with the specified configuration
func Initialize() func() {
	// Create logs directory in the application root
	logsDir := filepath.Join(".", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		os.Exit(1)
	}

	// Configure logging to both file and console
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		filepath.Join(logsDir, "godo.log"),
		"stdout",
	}
	config.ErrorOutputPaths = []string{
		filepath.Join(logsDir, "godo_error.log"),
		"stderr",
	}

	// Set log level based on environment
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if os.Getenv("DEBUG") == "1" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	// Create logger
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	log = logger
	return func() {
		_ = log.Sync()
	}
}

// InitializeFileOnly sets up logging to file only (no stdout)
func InitializeFileOnly() func() {
	// Create logger configuration
	cfg := zap.NewProductionConfig()

	// Set output paths to file only (no stdout)
	cfg.OutputPaths = []string{"logs/godo.log"}
	cfg.ErrorOutputPaths = []string{"logs/godo.log"}

	// Create logger
	zapLogger, err := cfg.Build()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Replace global logger
	zap.ReplaceGlobals(zapLogger)

	// Return cleanup function
	return func() {
		_ = zapLogger.Sync()
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	if log == nil {
		return nil
	}
	return log.Sync()
}

// Debug logs a debug message
func Debug(template string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Sugar().Debugf(template, args...)
}

// Info logs an info message
func Info(template string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Sugar().Infof(template, args...)
}

// Error logs an error message
func Error(template string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Sugar().Errorf(template, args...)
}

// Fatal logs a fatal message and exits
func Fatal(template string, args ...interface{}) {
	if log == nil {
		fmt.Printf(template+"\n", args...)
		os.Exit(1)
	}
	log.Sugar().Fatalf(template, args...)
}
