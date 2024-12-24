//go:build docker
// +build docker

package config

import (
	"golang.design/x/hotkey"
)

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
// This is a mock implementation for Docker environment
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModShift,
	}
}
