package logger

import "go.uber.org/zap"

// TestLogger returns a logger suitable for testing
func TestLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, _ := config.Build()
	return logger
}
