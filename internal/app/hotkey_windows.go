//go:build windows
// +build windows

package app

import "golang.design/x/hotkey"

type windowsHotkeyManager struct {
	hk *hotkey.Hotkey
}

func NewHotkeyManager() HotkeyManager {
	return &windowsHotkeyManager{}
}

// Implementation details...
