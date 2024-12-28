// Package quicknote implements the quick note window functionality
package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the quick note window
type Window struct {
	store  storage.Store
	logger logger.Logger
	window fyne.Window
}

// New creates a new quick note window
func New(store storage.Store, logger logger.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

// Show displays the quick note window
func (w *Window) Show() {
	w.window.Show()
}

// Hide hides the quick note window
func (w *Window) Hide() {
	w.window.Hide()
}
