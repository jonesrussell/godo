//go:build windows && !docker && !linux
// +build windows,!docker,!linux

package config

import (
	"strings"

	"golang.design/x/hotkey"
)

// Modifier constants for hotkeys on Windows
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
)

// GetDefaultQuickNoteKey returns the default key for quick note hotkey
func GetDefaultQuickNoteKey() string {
	return "N"
}

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		ModCtrl,
		ModShift,
	}
}

// GetDefaultQuickNoteKeyString returns the default key as a string
func GetDefaultQuickNoteKeyString() string {
	return GetDefaultQuickNoteKey()
}

// GetDefaultQuickNoteModifiersString returns the default modifiers as a string
func GetDefaultQuickNoteModifiersString() string {
	modifiers := []string{"Ctrl", "Shift"}
	return strings.Join(modifiers, "+")
}

// GetDefaultQuickNoteHotkey returns the complete default hotkey combination
func GetDefaultQuickNoteHotkey() HotkeyString {
	modifiers := []string{"Ctrl", "Shift"}
	return NewHotkeyString(modifiers, GetDefaultQuickNoteKey())
}
