package common

import (
	"fmt"
)

// LogConfig holds logging configuration
type LogConfig struct {
	Level       string   `yaml:"level"`
	Output      []string `yaml:"output"`
	ErrorOutput []string `yaml:"error_output"`
}

// HotkeyBinding represents a keyboard shortcut configuration
type HotkeyBinding struct {
	Modifiers []string `yaml:"modifiers"`
	Key       string   `yaml:"key"`
}

// String implements the Stringer interface for HotkeyBinding
func (h HotkeyBinding) String() string {
	return fmt.Sprintf("{%v %s}", h.Modifiers, h.Key)
}

// Error represents a domain error
type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a new domain error
func NewError(code string, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
