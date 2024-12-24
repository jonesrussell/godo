//go:build darwin
// +build darwin

package app

import "golang.design/x/hotkey"

type darwinHotkeyManager struct {
	hk *hotkey.Hotkey
}

func NewHotkeyManager() HotkeyManager {
	return &darwinHotkeyManager{}
}

// Implementation details...
