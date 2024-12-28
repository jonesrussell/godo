package mainwindow

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	store  storage.Store
	logger logger.Logger
	win    fyne.Window // Will be used in Setup implementation
}

// New creates a new main window
func New(store storage.Store, logger logger.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

// Setup initializes the window
func (w *Window) Setup() error {
	fmt.Println("Starting main window setup")
	w.logger.Debug("Starting main window setup")

	// Basic implementation to use the win field
	fmt.Println("Creating new window")
	w.logger.Debug("Creating new window")
	w.win = fyne.CurrentApp().NewWindow("Godo")
	if w.win == nil {
		fmt.Println("ERROR: Failed to create window - NewWindow returned nil")
		w.logger.Error("Failed to create window - NewWindow returned nil")
		return fmt.Errorf("failed to create window")
	}
	fmt.Println("Window created successfully")
	w.logger.Debug("Window created successfully")

	fmt.Println("Setting window properties")
	w.win.SetMaster() // Make this window the main window
	w.win.SetCloseIntercept(func() {
		fmt.Println("Window close requested")
		w.win.Hide()
	})

	fmt.Println("Resizing window")
	w.logger.Debug("Resizing window", "width", 800, "height", 600)
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.CenterOnScreen()

	fmt.Println("Showing window")
	w.logger.Debug("Showing window")
	w.win.Show()
	fmt.Println("Window setup complete")
	w.logger.Debug("Window setup complete")

	return nil
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.win
}
