// Package mainwindow implements the main application window
package mainwindow

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the main window functionality
type Window struct {
	store  storage.Store
	logger logger.Logger
	window fyne.Window
}

// New creates a new main window
func New(store storage.Store, logger logger.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

func (w *Window) Show() {
	w.window.Show()
}

func (w *Window) Hide() {
	w.window.Hide()
}

func (w *Window) SetContent(content fyne.CanvasObject) {
	w.window.SetContent(content)
}

func (w *Window) Resize(size fyne.Size) {
	w.window.Resize(size)
}

func (w *Window) CenterOnScreen() {
	w.window.CenterOnScreen()
}

func (w *Window) GetWindow() fyne.Window {
	return w.window
}
