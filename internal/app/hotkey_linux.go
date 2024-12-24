//go:build !docker
// +build !docker

package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"golang.design/x/hotkey"
)

// LinuxHotkeyManager implements HotkeyManager for Linux
type LinuxHotkeyManager struct {
	app *App
}

func NewLinuxHotkeyManager(app *App) HotkeyManager {
	return &LinuxHotkeyManager{app: app}
}

func (m *LinuxHotkeyManager) Setup() error {
	modifiers := config.GetDefaultQuickNoteModifiers()
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
