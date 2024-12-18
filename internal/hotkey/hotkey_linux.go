//go:build linux
// +build linux

package hotkey

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &defaultHotkeyManager{
			eventChan: make(chan struct{}),
		}, nil
	}
}
