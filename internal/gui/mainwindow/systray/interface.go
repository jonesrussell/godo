package systray

import "fyne.io/fyne/v2"

// Manager defines the behavior of a system tray icon
type Manager interface {
	SetupSystray(app fyne.App)
	Show()
	Hide()
}
