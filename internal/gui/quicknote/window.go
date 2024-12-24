package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
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
	logger *zap.Logger
	win    fyne.Window
}

// New creates a new quick note window
func New(store storage.Store, logger *zap.Logger) Interface {
	return newWindow(store, logger)
}

func newWindow(store storage.Store, logger *zap.Logger) Interface {
	return &Window{
		store:  store,
		logger: logger,
	}
}

func (w *Window) Setup() error {
	// Implementation
	return nil
}

func (w *Window) Hide() {
	if w.win != nil {
		w.win.Hide()
	}
}

func (w *Window) Show() {
	if w.win != nil {
		w.win.Show()
	}
}
