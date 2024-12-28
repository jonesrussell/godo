// Package mainwindow implements the main application window
package mainwindow

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the main window functionality
type Window struct {
	store  storage.TaskStore
	logger logger.Logger
	window fyne.Window
}

// New creates a new main window
func New(store storage.TaskStore, logger logger.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

// Show displays the window
func (w *Window) Show() {
	w.window.Show()
}

// Hide hides the window
func (w *Window) Hide() {
	w.window.Hide()
}

// SetContent sets the window's content
func (w *Window) SetContent(content fyne.CanvasObject) {
	w.window.SetContent(content)
}

// Resize changes the window's size
func (w *Window) Resize(size fyne.Size) {
	w.window.Resize(size)
}

// CenterOnScreen centers the window on the screen
func (w *Window) CenterOnScreen() {
	w.window.CenterOnScreen()
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.window
}
