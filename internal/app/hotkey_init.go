//go:build !docker

package app

func initHotkeyManager(app *App) HotkeyManager {
	return NewDefaultHotkeyManager(app)
}
