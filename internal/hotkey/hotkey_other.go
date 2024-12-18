//go:build !windows
// +build !windows

package hotkey

import (
	"context"

	hook "github.com/robotn/gohook"
)

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	eventChan chan struct{}
}

// NewHotkeyManager creates a new instance of HotkeyManager
func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{
		eventChan: make(chan struct{}, 1),
	}
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	go func() {
		hook.Register(hook.KeyDown, []string{"ctrl", "alt", "g"}, func(e hook.Event) {
			select {
			case h.eventChan <- struct{}{}:
			default:
				// Channel is full, skip this event
			}
		})
		s := hook.Start()
		<-ctx.Done()
		hook.End()
		<-s
	}()
	return nil
}

// GetEventChannel returns the channel that emits hotkey events
func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

// Cleanup performs any necessary cleanup
func (h *HotkeyManager) Cleanup() error {
	return nil
}
