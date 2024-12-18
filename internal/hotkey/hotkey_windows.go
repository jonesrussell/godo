//go:build windows
// +build windows

package hotkey

import "context"

type windowsHotkeyManager struct {
	config    BaseHotkeyConfig
	eventChan chan struct{}
}

func newPlatformHotkeyManager() (HotkeyManager, error) {
	return &windowsHotkeyManager{
		config:    DefaultConfig,
		eventChan: make(chan struct{}, 1),
	}, nil
}

func (h *windowsHotkeyManager) Start(ctx context.Context) error {
	// TODO: Implement Windows-specific hotkey handling
	<-ctx.Done()
	return nil
}

func (h *windowsHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}
