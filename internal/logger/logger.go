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

type zapLogger struct {
	log *zap.Logger
}

var defaultLogger = &zapLogger{}

func Initialize(config *common.LogConfig) (Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      config.Output,
		ErrorOutputPaths: config.ErrorOutput,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:   "ts",
			LevelKey:  "level",
			NameKey:   "logger",
			CallerKey: "caller",

			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	defaultLogger.log = logger
	return defaultLogger, nil
}

// Implementation of Logger interface
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Fatalw(msg, keysAndValues...)
}

// Package-level functions use defaultLogger
func Debug(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	defaultLogger.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warn(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	defaultLogger.Error(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatal(msg, keysAndValues...)
}
