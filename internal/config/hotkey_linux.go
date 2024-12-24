//go:build linux
// +build linux

package config

import "golang.design/x/hotkey"

// Modifier constants for hotkeys on Linux
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
	ModAlt   = hotkey.ModCtrl  // Alt key - using Ctrl as fallback since ModAlt is not supported on Linux
)

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		ModCtrl,
		ModShift, // Using Shift instead of Alt on Linux
	}
}
