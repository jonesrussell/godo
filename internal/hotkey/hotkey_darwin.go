//go:build darwin

package hotkey

import "context"

type darwinHotkeyManager struct {
	config    BaseHotkeyConfig
	eventChan chan struct{}
}

func newPlatformHotkeyManager() (HotkeyManager, error) {
	return &darwinHotkeyManager{
		config:    DefaultConfig,
		eventChan: make(chan struct{}, 1),
	}, nil
}

func (h *darwinHotkeyManager) Start(ctx context.Context) error {
	// TODO: Implement macOS-specific hotkey handling
	<-ctx.Done()
	return nil
}

func (h *darwinHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}
