// Package common provides shared types and utilities used across the application
package common

import (
	"fmt"
	"strings"
	"time"
)

// Config represents the main application configuration
type Config struct {
	Hotkeys HotkeyConfig `yaml:"hotkeys"`
	HTTP    HTTPConfig   `yaml:"http"`
	Logger  LogConfig    `yaml:"logger"`
}

// HotkeyConfig represents hotkey configuration
type HotkeyConfig struct {
	QuickNote HotkeyBinding `yaml:"quick_note"`
}

// HotkeyBinding represents a hotkey binding configuration
type HotkeyBinding struct {
	Modifiers []string `yaml:"modifiers"`
	Key       string   `yaml:"key"`
}

// String implements the Stringer interface for HotkeyBinding
func (h HotkeyBinding) String() string {
	return strings.Join(append(h.Modifiers, h.Key), "+")
}

// HTTPConfig represents HTTP server configuration
type HTTPConfig struct {
	Port              int `yaml:"port"`
	ReadTimeout       int `yaml:"read_timeout"`
	WriteTimeout      int `yaml:"write_timeout"`
	ReadHeaderTimeout int `yaml:"read_header_timeout"`
	IdleTimeout       int `yaml:"idle_timeout"`
}

// LogConfig represents logger configuration
type LogConfig struct {
	Level       string   `yaml:"level"`
	Console     bool     `yaml:"console"`
	File        bool     `yaml:"file"`
	FilePath    string   `yaml:"file_path"`
	Output      []string `yaml:"output"`
	ErrorOutput []string `yaml:"error_output"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if err := c.Hotkeys.QuickNote.Validate(); err != nil {
		return fmt.Errorf("invalid hotkey configuration: %w", err)
	}
	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("invalid HTTP configuration: %w", err)
	}
	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("invalid logger configuration: %w", err)
	}
	return nil
}

// Validate validates the hotkey binding configuration
func (h *HotkeyBinding) Validate() error {
	if len(h.Modifiers) == 0 {
		return fmt.Errorf("at least one modifier is required")
	}
	if h.Key == "" {
		return fmt.Errorf("key is required")
	}
	for _, mod := range h.Modifiers {
		switch mod {
		case "Ctrl", "Alt", "Shift":
			continue
		default:
			return fmt.Errorf("invalid modifier: %s", mod)
		}
	}
	return nil
}

// Validate validates the HTTP configuration
func (h *HTTPConfig) Validate() error {
	if h.Port <= 0 {
		return fmt.Errorf("port must be greater than 0")
	}
	if h.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be greater than 0")
	}
	if h.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be greater than 0")
	}
	if h.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("read header timeout must be greater than 0")
	}
	if h.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be greater than 0")
	}
	return nil
}

// Validate validates the logger configuration
func (l *LogConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[l.Level] {
		return fmt.Errorf("invalid log level: %s", l.Level)
	}
	if len(l.Output) == 0 {
		return fmt.Errorf("at least one output is required")
	}
	if len(l.ErrorOutput) == 0 {
		return fmt.Errorf("at least one error output is required")
	}
	if l.File && l.FilePath == "" {
		return fmt.Errorf("file path is required when file logging is enabled")
	}
	return nil
}

// GetReadTimeout returns the read timeout as time.Duration
func (h *HTTPConfig) GetReadTimeout() time.Duration {
	return time.Duration(h.ReadTimeout) * time.Second
}

// GetWriteTimeout returns the write timeout as time.Duration
func (h *HTTPConfig) GetWriteTimeout() time.Duration {
	return time.Duration(h.WriteTimeout) * time.Second
}

// GetReadHeaderTimeout returns the read header timeout as time.Duration
func (h *HTTPConfig) GetReadHeaderTimeout() time.Duration {
	return time.Duration(h.ReadHeaderTimeout) * time.Second
}

// GetIdleTimeout returns the idle timeout as time.Duration
func (h *HTTPConfig) GetIdleTimeout() time.Duration {
	return time.Duration(h.IdleTimeout) * time.Second
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
