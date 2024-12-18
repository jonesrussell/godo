package types

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
