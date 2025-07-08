// Package theme provides UI theme and asset management for the application
package theme

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed favicon.ico
var systrayIconBytes []byte

//go:embed icon.png
var appIconBytes []byte

// GetSystrayIconResource returns the system tray icon resource
func GetSystrayIconResource() fyne.Resource {
	return fyne.NewStaticResource("favicon.ico", systrayIconBytes)
}

// GetAppIconResource returns a Fyne resource for the application icon
func GetAppIconResource() fyne.Resource {
	return fyne.NewStaticResource("icon.png", appIconBytes)
}
