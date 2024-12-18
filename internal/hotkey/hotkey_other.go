//go:build !windows
// +build !windows

package hotkey

import (
	"context"
	"fmt"
	"runtime"

	"github.com/robotn/gohook"
)

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	eventChan chan struct{}
}

// NewHotkeyManager creates a new instance of HotkeyManager
func NewHotkeyManager() (*HotkeyManager, error) {
	// Check for supported platforms
	switch runtime.GOOS {
	case "darwin", "linux":
		// These platforms are supported
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return &HotkeyManager{
		eventChan: make(chan struct{}, 1),
	}, nil
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	// Add platform-specific key combinations
	keyCombo := []string{"ctrl", "alt", "g"}
	if runtime.GOOS == "darwin" {
		keyCombo = []string{"cmd", "alt", "g"} // Use cmd instead of ctrl on macOS
	}

	go func() {
		gohook.Register(gohook.KeyDown, keyCombo, func(e gohook.Event) {
			select {
			case h.eventChan <- struct{}{}:
			default:
				// Channel is full, skip this event
			}
		})
		s := gohook.Start()
		<-ctx.Done()
		gohook.End()
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
