//go:build windows && !ci && !android && !ios && !wasm && !test_web_driver

package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func createWindowsApp() fyne.App {
	return app.New()
}
