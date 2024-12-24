//go:build !docker
// +build !docker

package app

import (
	"fmt"

	"golang.design/x/hotkey"
)

// defaultHotkeyManager is the default implementation for non-Docker environments
type defaultHotkeyManager struct {
	app *App
	hk  *hotkey.Hotkey
}

// NewDefaultHotkeyManager creates a new default hotkey manager
func NewDefaultHotkeyManager(app *App) HotkeyManager {
	return &defaultHotkeyManager{app: app}
}

// Setup implements HotkeyManager interface
func (m *defaultHotkeyManager) Setup() error {
	// Unregister any existing hotkey
	if m.hk != nil {
		if err := m.hk.Unregister(); err != nil {
			return fmt.Errorf("failed to unregister existing hotkey: %w", err)
		}
	}

	// Register global hotkey
	m.hk = hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG)
	if err := m.hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	// Start hotkey listener
	go func() {
		for range m.hk.Keydown() {
			if m.app.quickNote != nil {
				m.app.quickNote.Show()
			}
		}
	}()

	return nil
}
