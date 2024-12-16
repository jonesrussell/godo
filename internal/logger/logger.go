package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// Open log file
	logFile, err := os.OpenFile(filepath.Join("logs", "godo.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// Create core that writes to both file and console
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(logFile),
			zapcore.DebugLevel,
		),
	)

	// Create logger with development options
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Development(),
	)
	log = logger.Sugar()
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// WithFields adds structured context to the logger
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	return log.With(fields)
}

// Sync flushes any buffered log entries
func Sync() error {
	// Ignore sync errors on stdout/stderr
	_ = log.Sync()
	return nil
}
