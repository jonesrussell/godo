// Package hotkey provides global hotkey functionality for the application
package hotkey

import (
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
)

// New creates a new platform-specific hotkey manager
func New(quickNote gui.QuickNoteManager, binding *common.HotkeyBinding, log logger.Logger) (Manager, error) {
	// Create Windows-specific manager
	manager, err := NewWindowsManager(log)
	if err != nil {
		return nil, err
	}

	// Set the quick note service and binding
	manager.SetQuickNote(quickNote, binding)

	return manager, nil
}
