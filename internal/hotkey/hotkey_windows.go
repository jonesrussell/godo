//go:build windows
// +build windows

package hotkey

import (
	"context"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/jonesrussell/godo/internal/logger"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procGetMessage       = user32.NewProc("GetMessageW")
)

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	eventChan chan struct{}
	config    HotkeyConfig
}

// NewHotkeyManager creates a new instance of HotkeyManager
func NewHotkeyManager() (*HotkeyManager, error) {
	manager := &HotkeyManager{
		eventChan: make(chan struct{}, 1),
		config:    DefaultConfig,
	}
	return manager, nil
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	ret, _, err := procRegisterHotKey.Call(
		uintptr(h.config.WindowHandle),
		uintptr(h.config.ID),
		uintptr(h.config.Modifiers),
		uintptr(h.config.Key),
	)

	if ret == 0 {
		return fmt.Errorf("failed to register hotkey: %v", err)
	}

	logger.Info("Successfully registered hotkey (ID=%d, Key='%c', Mods=0x%X)",
		h.config.ID, h.config.Key, h.config.Modifiers)

	// Start message loop in a goroutine
	go func() {
		var msg MSG

		for {
			select {
			case <-ctx.Done():
				// Unregister hotkey when context is cancelled
				ret, _, _ := procUnregisterHotKey.Call(
					uintptr(h.config.WindowHandle),
					uintptr(h.config.ID),
				)
				if ret == 0 {
					logger.Error("Failed to unregister hotkey")
				}
				return
			default:
				// GetMessage blocks until a message is received
				if ret, _, _ := procGetMessage.Call(
					uintptr(unsafe.Pointer(&msg)),
					0,
					0,
					0,
				); ret == 0 {
					// WM_QUIT received
					return
				}

				if msg.Message == WM_HOTKEY && msg.WParam == uintptr(h.config.ID) {
					select {
					case h.eventChan <- struct{}{}:
					default:
						// Channel is full, skip this event
					}
				}
			}
		}
	}()

	return nil
}

// GetEventChannel returns the channel that emits hotkey events
func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

// Cleanup performs any necessary cleanup
func (h *HotkeyManager) Cleanup() error {
	ret, _, err := procUnregisterHotKey.Call(
		uintptr(h.config.WindowHandle),
		uintptr(h.config.ID),
	)
	if ret == 0 {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}
	return nil
}
