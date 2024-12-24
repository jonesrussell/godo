//go:build docker && linux

package config

import (
	"golang.design/x/hotkey"
)

// GetDefaultQuickNoteModifiers returns empty modifiers for Docker environment
// since we don't actually register hotkeys in Docker
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{}
}
