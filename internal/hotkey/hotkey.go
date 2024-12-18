package hotkey

import "context"

// HotkeyManager handles global hotkey registration and events
type HotkeyManager interface {
	Start(ctx context.Context) error
	GetEventChannel() <-chan struct{}
	Cleanup() error
}

// NewHotkeyManager creates a platform-specific HotkeyManager instance
func NewHotkeyManager() (HotkeyManager, error) {
	return newPlatformHotkeyManager()
}
