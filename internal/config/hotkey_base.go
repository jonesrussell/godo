package config

import (
	"strings"

	"golang.design/x/hotkey"
)

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
}

// HotkeyFactory creates platform-specific hotkey handlers
type HotkeyFactory interface {
	NewHotkey([]hotkey.Modifier, hotkey.Key) HotkeyHandler
}

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote HotkeyString `mapstructure:"quick_note"`
}

// HotkeyString represents a hotkey combination as a string
type HotkeyString string

// NewHotkeyString creates a new hotkey string from modifiers and key
func NewHotkeyString(modifiers []string, key string) HotkeyString {
	return HotkeyString(strings.Join(append(modifiers, key), "+"))
}

// String implements the Stringer interface
func (h HotkeyString) String() string {
	return string(h)
}
