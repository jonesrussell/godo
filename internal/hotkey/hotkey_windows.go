//go:build windows
// +build windows

package hotkey

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &defaultHotkeyManager{
			eventChan: make(chan struct{}),
		}, nil
	}
}
