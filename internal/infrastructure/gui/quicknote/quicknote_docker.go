//go:build docker

package quicknote

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// dockerWindow represents a no-op quick note window for Docker environments
type dockerWindow struct {
	store storage.Store
	log   logger.Logger
}

// newWindow creates a new quick note window for Docker environments
func newWindow(store storage.Store) Interface {
	return &dockerWindow{
		store: store,
	}
}

// Initialize sets up the window with the given app and logger
func (w *dockerWindow) Initialize(app fyne.App, log logger.Logger) {
	w.log = log
}

// Show displays the quick note window (no-op in Docker)
func (w *dockerWindow) Show() {
	// No-op in Docker environment
}

// Hide hides the quick note window (no-op in Docker)
func (w *dockerWindow) Hide() {
	// No-op in Docker environment
}
