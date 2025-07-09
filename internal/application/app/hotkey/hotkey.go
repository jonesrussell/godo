// Package hotkey provides global hotkey functionality for the application
package hotkey

import (
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

// New creates a new platform-specific hotkey manager
func New(quickNote QuickNoteService, binding *config.HotkeyBinding, log logger.Logger) (Manager, error) {
	// Create platform-specific manager
	manager := newPlatformManager(quickNote, binding)

	// Set logger if the manager supports it
	if logSetter, ok := manager.(interface{ SetLogger(logger.Logger) }); ok {
		logSetter.SetLogger(log)
	}

	return manager, nil
}
