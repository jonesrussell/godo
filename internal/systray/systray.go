package systray

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

// Manager defines the interface for platform-specific systray implementations
type Manager interface {
	Setup() error
	AddMenuItem(label, tooltip string) *systray.MenuItem
	Quit()
}

// SetupSystray is a convenience function for setting up the systray
func SetupSystray() (Manager, error) {
	manager := newManager()
	if err := manager.Setup(); err != nil {
		return nil, err
	}
	return manager, nil
}

// Variable to hold the platform-specific systray manager constructor
var newManager = func() Manager {
	return &defaultManager{}
}

// defaultManager provides a fallback implementation
type defaultManager struct {
	ready bool
}

func (s *defaultManager) Setup() error {
	icon, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load icon: %v", err)
		return fmt.Errorf("failed to load icon: %w", err)
	}

	go systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTooltip("Godo - Quick Note Todo App")
		s.ready = true
	}, func() {
		// Cleanup on exit
	})

	return nil
}

func (s *defaultManager) AddMenuItem(label, tooltip string) *systray.MenuItem {
	return systray.AddMenuItem(label, tooltip)
}

func (s *defaultManager) Quit() {
	systray.Quit()
}
