// Package hotkey provides hotkey management functionality for the application
package hotkey

import (
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
)

// Manager defines the interface for hotkey management
type Manager interface {
	// Register registers a hotkey with the system
	Register(binding *common.HotkeyBinding) error
	// Unregister removes the hotkey registration from the system
	Unregister() error
	// Start begins listening for hotkey events
	Start() error
	// Stop ends the hotkey listening
	Stop() error
	// IsRegistered returns whether a hotkey is currently registered
	IsRegistered() bool
}

// NewManager creates a new hotkey manager based on the platform
func NewManager(log logger.Logger) (Manager, error) {
	handler := NewWindowsHotkeyHandler()
	return NewWindowsHotkeyManager(log, handler), nil
}
