package ui

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/logger"
)

// SetupSystray initializes the system tray icon and menu
func SetupSystray() error {
	icon, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load icon: %v", err)
		return fmt.Errorf("failed to load icon: %w", err)
	}

	systray.SetIcon(icon)
	systray.SetTooltip("Godo - Quick Note Todo App")

	// Platform-specific title setting is handled in the platform-specific files
	setPlatformSpecificTitle()

	return nil
}

// Platform-specific implementations will override this
var setPlatformSpecificTitle = func() {
	// Default implementation does nothing
}
