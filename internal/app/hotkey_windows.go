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

func (m *windowsHotkeyManager) Register() error {
	// Implementation
	return nil
}

func (m *windowsHotkeyManager) Unregister() error {
	// Implementation
	return nil
}
