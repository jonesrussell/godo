//go:build windows && !linux

package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func initApp() fyne.App {
	return app.New()
}
