//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"github.com/jonesrussell/godo/internal/common"
)

// Handler defines the interface for platform-specific hotkey handling
type Handler interface {
	// Register registers the hotkey with the system
	Register(*common.HotkeyBinding) error

	// Unregister removes the hotkey registration from the system
	Unregister(*common.HotkeyBinding) error

	// Start begins listening for hotkey events
	Start() error

	// Stop ends the hotkey listening
	Stop() error

	// IsRegistered returns whether the hotkey is currently registered
	IsRegistered() bool
}
