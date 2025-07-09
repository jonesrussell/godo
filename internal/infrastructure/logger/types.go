// Package logger provides logging functionality
package logger

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string   `mapstructure:"level"`
	Console  bool     `mapstructure:"console"`
	File     bool     `mapstructure:"file"`
	FilePath string   `mapstructure:"file_path"`
	Output   []string `mapstructure:"output"`
}

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}
