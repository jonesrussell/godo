package app

import "fyne.io/fyne/v2"

// UI defines the user interface operations
type UI interface {
	Show()
	Hide()
	SetContent(content fyne.CanvasObject)
	Resize(size fyne.Size)
	CenterOnScreen()
}

// Application defines the core application behavior
type Application interface {
	SetupUI()
	Run()
	Cleanup()
}

// QuickNoteService defines quick note operations
type QuickNoteService interface {
	Show()
	Hide()
}

// SystemTrayService defines system tray operations
type SystemTrayService interface {
	Setup(menu *fyne.Menu)
	SetIcon(resource fyne.Resource)
}
