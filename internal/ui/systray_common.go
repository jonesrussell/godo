package ui

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

// SystrayManager defines the interface for platform-specific systray implementations
type SystrayManager interface {
	Setup() error
}

// defaultSystray provides a default implementation of SystrayManager
type defaultSystray struct{}

// Setup initializes the system tray icon and menu
func (s *defaultSystray) Setup() error {
	icon, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load icon: %v", err)
		return fmt.Errorf("failed to load icon: %w", err)
	}

	go systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTooltip("Godo - Quick Note Todo App")
	}, func() {
		// Cleanup on exit
	})

	return nil
}

// SetupSystray is a convenience function for setting up the systray
func SetupSystray() error {
	manager := newSystrayManager()
	return manager.Setup()
}

// Variable to hold the platform-specific systray manager constructor
var newSystrayManager = func() SystrayManager {
	return &defaultSystray{}
}
