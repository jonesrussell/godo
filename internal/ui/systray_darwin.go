//go:build darwin

package ui

import "github.com/getlantern/systray"

// SetupSystray initializes the system tray icon and menu for macOS
func SetupSystray(icon []byte) {
	systray.SetIcon(icon)
	systray.SetTitle("Godo")
	systray.SetTooltip("Godo")
}
