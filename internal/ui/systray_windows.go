//go:build windows

package ui

import "github.com/getlantern/systray"

// SetupSystray initializes the system tray icon and menu for Windows
func SetupSystray(icon []byte) {
	systray.SetIcon(icon)
	// Windows doesn't support SetTitle
	systray.SetTooltip("Godo")
}
