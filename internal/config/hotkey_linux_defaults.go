//go:build linux && !windows && !darwin && !docker
// +build linux,!windows,!darwin,!docker

package config

import (
	"strings"

	"golang.design/x/hotkey"
)

// Modifier constants for hotkeys on Linux
const (
	ModCtrl  = hotkey.ModCtrl  // Control key
	ModShift = hotkey.ModShift // Shift key
	ModAlt   = hotkey.Mod1     // Alt key (Mod1 on X11)
)

// GetDefaultQuickNoteKey returns the default key for quick note hotkey
func GetDefaultQuickNoteKey() string {
	return "N"
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
