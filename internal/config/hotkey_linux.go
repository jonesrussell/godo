//go:build !windows && !android && !ios && !wasm && !js

package config

import (
	"golang.design/x/hotkey"
)

// GetDefaultQuickNoteModifiers returns the default modifiers for quick note hotkey
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModShift,
		hotkey.Mod2, // Typically maps to Shift on X11
		hotkey.Mod4, // Typically maps to Super/Windows key on X11
	}
}
