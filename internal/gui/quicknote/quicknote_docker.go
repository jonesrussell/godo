//go:build docker
// +build docker

package quicknote

import (
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Config holds the configuration for the quick note window
type Config struct {
	Store  storage.Store
	Logger logger.Logger
}

// Window represents the quick note window
type Window struct {
	config Config
}

// New creates a new quick note window
func New(store storage.Store) *Window {
	return &Window{
		config: Config{
			Store: store,
		},
	}
}

// Show displays the quick note window
func (w *Window) Show() {
	// No-op in Docker environment
}

// Hide hides the quick note window
func (w *Window) Hide() {
	// No-op in Docker environment
}
