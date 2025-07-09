// Package theme provides UI theme and asset management for the application
package theme

import (
	"fyne.io/fyne/v2"

	_ "embed"
)

//go:embed icon.ico
var iconData []byte

// AppIcon returns the application icon resource
// This is used for both the main application and system tray
func AppIcon() fyne.Resource {
	return fyne.NewStaticResource("icon.ico", iconData)
}
