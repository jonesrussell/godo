//go:build darwin
// +build darwin

package ui

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

func init() {
	newSystrayManager = func() SystrayManager {
		return &darwinSystray{}
	}
}

type darwinSystray struct{}

func (s *darwinSystray) Setup() error {
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
