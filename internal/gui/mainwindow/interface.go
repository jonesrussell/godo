// Package mainwindow implements the main application window
package mainwindow

import "fyne.io/fyne/v2"

// Interface defines the main window functionality
type Interface interface {
	Show()
	Hide()
	SetContent(content fyne.CanvasObject)
	Resize(size fyne.Size)
	CenterOnScreen()
	GetWindow() fyne.Window
	Refresh()
}
