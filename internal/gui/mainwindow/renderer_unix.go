//go:build !windows && !android && !ios && !wasm && !js

package mainwindow

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func createUnixApp() fyne.App {
	// Use software renderer if specified
	if os.Getenv("FYNE_RENDERER") == "software" {
		return app.NewWithDriver(&app.DriverConfig{
			Renderer: app.RendererSoftware,
		})
	}
	return app.New()
}
