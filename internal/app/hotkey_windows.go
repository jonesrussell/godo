//go:build windows
// +build windows

package app

import (
	"golang.design/x/hotkey"
)

func getHotkeyModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModAlt,
	}
}
