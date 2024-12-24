//go:build linux
// +build linux

package app

import "golang.design/x/hotkey"

type linuxHotkeyManager struct {
	hk *hotkey.Hotkey
}

func NewHotkeyManager() HotkeyManager {
	return &linuxHotkeyManager{}
}

// Implementation details...
