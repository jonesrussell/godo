package config

import (
	"strings"
)

// Modifier represents a hotkey modifier (Ctrl, Alt, etc.)
type Modifier uint8

// Key represents a keyboard key
type Key uint16

// Event represents a hotkey event
type Event struct{}

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register() error
	Unregister() error
	Keydown() <-chan Event
}

// HotkeyFactory creates platform-specific hotkey handlers
type HotkeyFactory interface {
	NewHotkey([]Modifier, Key) HotkeyHandler
}

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote HotkeyCombo `mapstructure:"quick_note"`
}

// HotkeyCombo represents a hotkey combination as a string
type HotkeyCombo string

// NewHotkeyCombo creates a new hotkey string from modifiers and key
func NewHotkeyCombo(modifiers []string, key string) HotkeyCombo {
	return HotkeyCombo(strings.Join(append(modifiers, key), "+"))
}

// String returns the string representation of the hotkey combination
func (h HotkeyCombo) String() string {
	return string(h)
}
