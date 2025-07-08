// Package hotkey provides global hotkey functionality for the application
package hotkey

import (
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/shared/common"
)

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}

// New creates a new platform-specific hotkey manager
func New(quickNote QuickNoteService, binding *common.HotkeyBinding, log logger.Logger) (Manager, error) {
	// Create platform-specific manager
	manager := newPlatformManager(quickNote, binding)

	// Set logger if the manager supports it
	if logSetter, ok := manager.(interface{ SetLogger(logger.Logger) }); ok {
		logSetter.SetLogger(log)
	}

	return manager, nil
}
