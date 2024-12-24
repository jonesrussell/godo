//go:build linux
// +build linux

package app

import (
	"golang.design/x/hotkey"
)

func getHotkeyModifiers() []hotkey.Modifier {
	return []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.Mod1,
	}
}
