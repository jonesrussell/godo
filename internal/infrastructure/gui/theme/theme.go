// Package theme provides UI theme and asset management for the application
package theme

import (
	"runtime"

	"fyne.io/fyne/v2"

	_ "embed"
)

//go:embed icon.png
var iconData []byte

//go:embed favicon.ico
var iconICOData []byte

// AppIcon returns the application icon resource
// This is used for both the main application and system tray
func AppIcon() fyne.Resource {
	return fyne.NewStaticResource("icon.png", iconData)
}

// SystrayIcon returns the appropriate icon for the system tray
// Windows typically expects ICO format for systray icons
func SystrayIcon() fyne.Resource {
	if runtime.GOOS == "windows" {
		return fyne.NewStaticResource("favicon.ico", iconICOData)
	}
	return AppIcon()
}

// Deprecated: Use AppIcon() instead
// GetAppIconResource returns the application icon resource
func GetAppIconResource() fyne.Resource {
	return AppIcon()
}

// Deprecated: Use SystrayIcon() instead
// GetSystrayIconResource returns the system tray icon resource
func GetSystrayIconResource() fyne.Resource {
	return SystrayIcon()
}
