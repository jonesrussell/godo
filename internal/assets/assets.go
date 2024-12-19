package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"

	"github.com/jonesrussell/godo/internal/logger"
)

//go:embed favicon.ico
var systrayIconBytes []byte

//go:embed Icon.png
var appIconBytes []byte

// GetSystrayIconResource returns a Fyne resource for the system tray icon
func GetSystrayIconResource() fyne.Resource {
	logger.Debug("Loading system tray icon")
	return fyne.NewStaticResource("systray", systrayIconBytes)
}

// GetAppIconResource returns a Fyne resource for the application icon
func GetAppIconResource() fyne.Resource {
	logger.Debug("Loading application icon")
	return fyne.NewStaticResource("app", appIconBytes)
}
