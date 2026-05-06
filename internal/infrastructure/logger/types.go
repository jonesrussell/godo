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
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
	WithError(err error) Logger
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}
