// Package common provides shared types and utilities used across the application
package common

import "time"

// LogConfig holds logging configuration
type LogConfig struct {
	Level       string   `mapstructure:"level" yaml:"level"`
	Console     bool     `mapstructure:"console" yaml:"console"`
	File        bool     `mapstructure:"file" yaml:"file"`
	FilePath    string   `mapstructure:"file_path" yaml:"file_path"`
	Output      []string `mapstructure:"output" yaml:"output"`
	ErrorOutput []string `mapstructure:"error_output" yaml:"error_output"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port              int `mapstructure:"port" yaml:"port"`
	ReadTimeout       int `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout      int `mapstructure:"write_timeout" yaml:"write_timeout"`
	ReadHeaderTimeout int `mapstructure:"read_header_timeout" yaml:"read_header_timeout"`
	IdleTimeout       int `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

// GetReadTimeout returns the read timeout as a duration
func (c *HTTPConfig) GetReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeout) * time.Second
}

// GetWriteTimeout returns the write timeout as a duration
func (c *HTTPConfig) GetWriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeout) * time.Second
}

// GetReadHeaderTimeout returns the read header timeout as a duration
func (c *HTTPConfig) GetReadHeaderTimeout() time.Duration {
	return time.Duration(c.ReadHeaderTimeout) * time.Second
}

// GetIdleTimeout returns the idle timeout as a duration
func (c *HTTPConfig) GetIdleTimeout() time.Duration {
	return time.Duration(c.IdleTimeout) * time.Second
}

// HotkeyBinding represents a keyboard shortcut configuration
type HotkeyBinding struct {
	Modifiers []string `yaml:"modifiers"`
	Key       string   `yaml:"key"`
}

// String implements the Stringer interface for HotkeyBinding
func (h HotkeyBinding) String() string {
	return h.Key
}

// Error represents a domain error
type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Code + ": " + e.Message + ": " + e.Err.Error()
	}
	return e.Code + ": " + e.Message
}

// NewError creates a new domain error
func NewError(code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
