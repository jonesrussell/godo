//go:build linux && !windows

package mainwindow

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func initApp() fyne.App {
	// Use software renderer for cross-compilation
	if os.Getenv("FYNE_RENDERER") == "software" {
		return app.NewWithDriver(app.NewSoftwareDriver())
	}
	return app.New()
}
