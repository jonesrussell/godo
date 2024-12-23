package logger

import "go.uber.org/zap"

// NewZapLogger creates a new Logger from a zap.Logger instance
func NewZapLogger(z *zap.Logger) Logger {
	return &zapLogger{z.Sugar()}
}

// Implement the interface methods
func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) WithError(err error) Logger {
	return &zapLogger{l.SugaredLogger.With("error", err)}
}

func (l *zapLogger) WithField(key string, value interface{}) Logger {
	return &zapLogger{l.SugaredLogger.With(key, value)}
}

func (l *zapLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{l.SugaredLogger.With(args...)}
}
