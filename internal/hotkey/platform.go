//go:build !windows && !darwin && !linux

package hotkey

// newPlatformHotkeyManager is implemented differently for each platform
// The implementation is in the platform-specific files:
// - hotkey_windows.go
// - hotkey_darwin.go
// - hotkey_linux.go
func newPlatformHotkeyManager() (HotkeyManager, error) {
	return nil, ErrUnsupportedPlatform
}
