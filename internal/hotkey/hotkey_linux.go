//go:build linux

package hotkey

import (
	"context"

	hook "github.com/robotn/gohook"
)

type linuxHotkeyManager struct {
	config    BaseHotkeyConfig
	eventChan chan struct{}
}

func newPlatformHotkeyManager() (HotkeyManager, error) {
	return &linuxHotkeyManager{
		config:    DefaultConfig,
		eventChan: make(chan struct{}, 1),
	}, nil
}

func (h *linuxHotkeyManager) Start(ctx context.Context) error {
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		// Check if the hotkey combination matches
		if e.Keycode == int(h.config.Key) && e.Rawcode == h.config.Modifiers {
			select {
			case <-ctx.Done():
				return
			case h.eventChan <- struct{}{}:
				// Signal sent successfully
			default:
				// Channel is full, skip this event
			}
		}
	})

	go func() {
		<-ctx.Done()
		hook.End()
	}()

	go hook.Start()
	return nil
}

func (h *linuxHotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}
