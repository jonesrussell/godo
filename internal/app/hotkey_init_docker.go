//go:build docker && !ci && !android && !ios && !wasm && !test_web_driver

package app

func initHotkeyManager(app *App) HotkeyManager {
	return NewNoopHotkeyManager(app)
}
