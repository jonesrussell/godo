//go:build docker

package app

// noopHotkeyManager is a no-op implementation for Docker environments
type noopHotkeyManager struct {
	app *App
}

// NewNoopHotkeyManager creates a new no-op hotkey manager
func NewNoopHotkeyManager(app *App) HotkeyManager {
	return &noopHotkeyManager{app: app}
}

// Setup implements HotkeyManager interface
func (m *noopHotkeyManager) Setup() error {
	return nil
}
