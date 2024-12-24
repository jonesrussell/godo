//go:build windows || darwin
// +build windows darwin

package config

import "golang.design/x/hotkey"

// Modifier constants for hotkeys on Windows and Darwin
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
	ModAlt   = hotkey.ModAlt   // Alt/Option key
)

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		ModCtrl,
		ModAlt,
	}
}
