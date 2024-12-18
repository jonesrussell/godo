package hotkey

import "context"

// HotkeyManager defines the interface for platform-specific hotkey implementations
type HotkeyManager interface {
	// Start begins listening for hotkey events
	Start(ctx context.Context) error
	// GetEventChannel returns a channel that receives events when the hotkey is pressed
	GetEventChannel() <-chan struct{}
	// Stop stops listening for hotkey events and cleans up resources
	Stop() error
}

// NewHotkeyManager creates a new platform-specific hotkey manager
func NewHotkeyManager() (HotkeyManager, error) {
	return newPlatformHotkeyManager()
}
