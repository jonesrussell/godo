//go:build windows || darwin
// +build windows darwin

package config

import (
	"strings"

	"golang.design/x/hotkey"
)

// Modifier constants for hotkeys on Windows and Darwin
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
	ModAlt   = hotkey.ModAlt   // Alt/Option key
)

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote HotkeyString `mapstructure:"quick_note"`
}

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		ModCtrl,
		ModAlt,
	}
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
