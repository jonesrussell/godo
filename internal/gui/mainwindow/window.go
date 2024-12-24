package mainwindow

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
)

// Window represents the main application window
type Window struct {
	store  storage.Store
	logger *zap.Logger
	win    fyne.Window // Will be used in Setup implementation
}

// New creates a new main window
func New(store storage.Store, logger *zap.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

// Setup initializes the window
func (w *Window) Setup() error {
	// Basic implementation to use the win field
	w.win = fyne.CurrentApp().NewWindow("Godo")
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.Show()
	return nil
}
