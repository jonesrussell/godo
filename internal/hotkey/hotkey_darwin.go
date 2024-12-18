//go:build darwin
// +build darwin

package hotkey

import "context"

func init() {
	newPlatformHotkeyManager = func() (HotkeyManager, error) {
		return &darwinHotkeyManager{
			eventChan: make(chan struct{}),
		}, nil
	}
}

type darwinHotkeyManager struct {
	eventChan chan struct{}
}

func (h *darwinHotkeyManager) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (h *darwinHotkeyManager) Stop() error {
	return nil
}

func (h *darwinHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}
