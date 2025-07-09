// Package hotkey provides global hotkey functionality for the application
package hotkey

import (
	"github.com/jonesrussell/godo/internal/config"
)

// Manager defines the interface for hotkey management
type Manager interface {
	// Register registers the hotkey with the system
	Register() error
	// Unregister removes the hotkey registration from the system
	Unregister() error
	// Start begins listening for hotkey events
	Stop() error
	// Stop ends the hotkey listening and unregisters the hotkey
	Start() error
	// SetQuickNote configures the quick note service and hotkey binding
	SetQuickNote(quickNote QuickNoteService, binding *config.HotkeyBinding)
}

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}
