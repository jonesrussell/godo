//go:build !docker
// +build !docker

package app

import "golang.design/x/hotkey"

// defaultHotkeyManager is the default implementation for non-Docker environments
type defaultHotkeyManager struct {
	app *App
}

// NewDefaultHotkeyManager creates a new default hotkey manager
func NewDefaultHotkeyManager(app *App) HotkeyManager {
	return &defaultHotkeyManager{app: app}
}

// Setup implements HotkeyManager interface
func (m *defaultHotkeyManager) Setup() error {
	// Register global hotkey
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyN)
	if err := hk.Register(); err != nil {
		return err
	}

	// Start hotkey listener
	go func() {
		for range hk.Keydown() {
			if m.app.quickNote != nil {
				m.app.quickNote.Show()
			}
		}
	}()

	return nil
}
