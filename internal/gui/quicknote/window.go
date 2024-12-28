package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Interface defines the behavior of a quick note window
type Interface interface {
	Setup() error
	Hide()
	Show()
}

// Window implements the quick note window
type Window struct {
	store  storage.Store
	logger logger.Logger
	win    fyne.Window
}

// New creates a new quick note window
func New(store storage.Store, logger logger.Logger) Interface {
	return newWindow(store, logger)
}

func newWindow(store storage.Store, logger logger.Logger) Interface {
	return &Window{
		store:  store,
		logger: logger,
	}
}

func (w *Window) Setup() error {
	w.logger.Debug("Setting up quick note window")
	
	// Create the window
	w.win = fyne.CurrentApp().NewWindow("Quick Note")
	w.win.Resize(fyne.NewSize(400, 300))
	w.win.CenterOnScreen()
	
	// TODO: Add text input and save button
	
	w.logger.Debug("Quick note window setup complete")
	return nil
}

func (w *Window) Hide() {
	if w.win != nil {
		w.logger.Debug("Hiding quick note window")
		w.win.Hide()
	} else {
		w.logger.Error("Cannot hide quick note window - window is nil")
	}
}

func (w *Window) Show() {
	if w.win != nil {
		w.logger.Debug("Showing quick note window")
		w.win.Show()
		w.win.CenterOnScreen()
		w.win.RequestFocus()
	} else {
		w.logger.Error("Cannot show quick note window - window is nil")
	}
}
