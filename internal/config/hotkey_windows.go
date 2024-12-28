//go:build windows && !docker && !linux
// +build windows,!docker,!linux

package config

import (
	"golang.design/x/hotkey"
)

// Modifier constants for hotkeys on Windows
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
)

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		ModCtrl,
		ModShift,
	}
}
