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
	SetupUI() error
	Run()
	Cleanup()
}
