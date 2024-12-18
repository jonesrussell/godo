//go:build linux

package ui

import "github.com/getlantern/systray"

// SetupSystray initializes the system tray icon and menu for Linux
func SetupSystray(icon []byte) {
	systray.SetIcon(icon)
	// Linux typically uses AppIndicator/StatusNotifier
	systray.SetTooltip("Godo")
}
