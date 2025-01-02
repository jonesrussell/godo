//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

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
	return "G"
}

// GetDefaultQuickNoteKeyString returns the default key as a string
func GetDefaultQuickNoteKeyString() string {
	return strings.ToUpper(GetDefaultQuickNoteKey())
}

// GetDefaultQuickNoteModifiersString returns the default modifiers as a string
func GetDefaultQuickNoteModifiersString() string {
	return "Ctrl+Shift"
}

// GetDefaultQuickNoteHotkey returns the default hotkey configuration
func GetDefaultQuickNoteHotkey() HotkeyString {
	return HotkeyString{
		Key:       GetDefaultQuickNoteKeyString(),
		Modifiers: []string{"ctrl", "shift"},
	}
}

// GetDefaultQuickNoteCombo returns the default hotkey combination
func GetDefaultQuickNoteCombo() HotkeyCombo {
	return NewHotkeyCombo([]string{"Ctrl", "Shift"}, GetDefaultQuickNoteKey())
}
