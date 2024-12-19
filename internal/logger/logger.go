package logger

import (
	"github.com/jonesrussell/godo/internal/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
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
	log *zap.Logger
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Fatalw(msg, keysAndValues...)
}

// NewZapLogger creates a new ZapLogger instance
func NewZapLogger(config *Config) (Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Configure output paths
	if config.Console {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, "stdout")
		zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, "stderr")
	}

	if config.File && config.FilePath != "" {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, config.FilePath)
		zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, config.FilePath)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{log: logger}, nil
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
