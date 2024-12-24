package systray

import "fyne.io/fyne/v2"

// Interface defines the system tray functionality
type Interface interface {
	Setup(*fyne.Menu)
	SetIcon(fyne.Resource)
	IsReady() bool
}
