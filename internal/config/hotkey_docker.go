//go:build docker && !ci && !android && !ios && !wasm && !test_web_driver

package config

import (
	"golang.design/x/hotkey"
)

// GetDefaultQuickNoteModifiers returns empty modifiers for Docker environment
// since we don't actually register hotkeys in Docker
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{}
}
