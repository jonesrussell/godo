//go:build docker
// +build docker

package mainwindow

import (
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	store storage.Store
}

// New creates a new main window
func New(store storage.Store) *Window {
	return &Window{
		store: store,
	}
}

// Show displays the main window
func (w *Window) Show() {
	// No-op in Docker environment
}

// Setup initializes the window
func (w *Window) Setup() {
	// No-op in Docker environment
}
