package app

import (
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// TestApp is a special version of App used for testing
type TestApp struct {
	Logger     logger.Logger
	Store      storage.Store
	MainWindow gui.MainWindow
	QuickNote  gui.QuickNote
	Hotkey     hotkey.Manager
}

// SetupUI implements the Application interface for testing
func (a *TestApp) SetupUI() {}

// Run implements the Application interface for testing
func (a *TestApp) Run() {}

// Cleanup implements the Application interface for testing
func (a *TestApp) Cleanup() {}