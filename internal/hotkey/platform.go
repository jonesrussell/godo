//go:build !windows && !darwin && !linux

package hotkey

import "context"

// HotkeyManager defines the interface for platform-specific hotkey implementations
type HotkeyManager interface {
	Start(ctx context.Context) error
	Stop() error
	GetEventChannel() <-chan struct{}
}

// newPlatformHotkeyManager is implemented differently for each platform
func newPlatformHotkeyManager() (HotkeyManager, error) {
	return nil, ErrUnsupportedPlatform
}
