package logger

import (
	"errors"

	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	// Helper methods
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// Config holds logger configuration
type Config struct {
	Level    string
	Console  bool
	File     bool
	FilePath string
}

// ZapLogger is the concrete implementation of Logger using zap
type ZapLogger struct {
	log *zap.SugaredLogger
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	return &ZapLogger{
		log: l.log.With(key, value),
	}
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	// Convert map to key-value pairs
	kvs := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		kvs = append(kvs, k, v)
	}
	return &ZapLogger{
		log: l.log.With(kvs...),
	}
}

func (l *ZapLogger) WithError(err error) Logger {
	return &ZapLogger{
		log: l.log.With("error", err),
	}
}

// NewZapLogger creates a new ZapLogger instance
func NewZapLogger(cfg *Config) (Logger, error) {
	// Create the basic configuration
	zapCfg := zap.NewProductionConfig()

	// Set log level
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, errors.New("invalid log level")
	}
	zapCfg.Level = level

	// Configure output paths (avoid duplicates)
	zapCfg.OutputPaths = []string{"stdout"}      // Single output path
	zapCfg.ErrorOutputPaths = []string{"stderr"} // Single error output path

	// Build the logger
	logger, err := zapCfg.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, errors.New("failed to build logger")
	}

	return &ZapLogger{
		log: logger.Sugar(),
	}, nil
}

// New creates a new logger instance (for backward compatibility)
func New(config *common.LogConfig) (Logger, error) {
	return NewZapLogger(&Config{
		Level:    config.Level,
		Console:  true,
		File:     false,
		FilePath: "",
	})
}
