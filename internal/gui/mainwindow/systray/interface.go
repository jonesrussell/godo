package systray

import "fyne.io/fyne/v2"

// SystrayManager defines the behavior of a system tray icon
type SystrayManager interface {
	SetupSystray(app fyne.App)
	Show()
	Hide()
}
