//go:build linux

package config

import (
	"golang.design/x/hotkey"
)

// On Linux (X11), some keys may be mapped to multiple Mod keys.
// We need to use the correct underlying keycode combination.
// For example, Ctrl+Shift+G might be registered as: Ctrl+Mod2+Mod4+G
func GetDefaultQuickNoteModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModShift,
		hotkey.Mod2, // Typically maps to Shift on X11
		hotkey.Mod4, // Typically maps to Super/Windows key on X11
	}
}
