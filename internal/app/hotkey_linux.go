//go:build linux && !docker
// +build linux,!docker

package app

import (
	"golang.design/x/hotkey"
)

// LinuxHotkeyManager implements HotkeyManager for Linux
type LinuxHotkeyManager struct {
	app *App
}

// NewHotkeyManager creates a new hotkey manager for Linux
func NewHotkeyManager(app *App) HotkeyManager {
	return &LinuxHotkeyManager{app: app}
}

func (m *LinuxHotkeyManager) Setup() error {
	modifiers := []hotkey.Modifier{
		hotkey.ModCtrl,
		hotkey.ModShift,
		hotkey.Mod2, // Typically maps to Shift on X11
		hotkey.Mod4, // Typically maps to Super/Windows key on X11
	}
	key := hotkey.KeyN

	hk := hotkey.New(modifiers, key)
	if err := hk.Register(); err != nil {
		return err
	}

	go func() {
		for range hk.Keydown() {
			m.app.quickNote.Show()
		}
	}()

	return nil
}
