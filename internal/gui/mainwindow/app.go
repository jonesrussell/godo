//go:build !docker

package mainwindow

import "fyne.io/fyne/v2"

// createApp creates a new Fyne application with platform-specific configuration
func createApp() fyne.App {
	return initApp()
}
