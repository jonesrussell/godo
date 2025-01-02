// Package hotkey provides global hotkey functionality for the application
package hotkey

import (
	"fmt"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
)

// New creates a new hotkey manager with the given configuration
func New(quickNote gui.QuickNoteManager, binding *common.HotkeyBinding, log logger.Logger) (Manager, error) {
	// Create Windows-specific manager
	manager, err := NewWindowsManager(log)
	if err != nil {
		return nil, fmt.Errorf("failed to create hotkey manager: %w", err)
	}

	// Set the quick note service and binding
	manager.SetQuickNote(quickNote, binding)

	return manager, nil
}
