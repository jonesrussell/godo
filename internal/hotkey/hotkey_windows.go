//go:build windows
// +build windows

package hotkey

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
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
	}
	return manager, nil
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	// Cleanup any existing registration
	_, _ = unregisterHotkey(h.config.WindowHandle, h.config.ID)
	time.Sleep(100 * time.Millisecond)

	success, err := registerHotkey(h.config)
	if !success {
		lastErr := syscall.GetLastError()
		return fmt.Errorf("failed to register hotkey: %w (lastErr=%d)", err, lastErr)
	}

	logger.Info("Successfully registered hotkey (ID=%d, Key='%c', Mods=0x%X)",
		h.config.ID, h.config.Key, h.config.Modifiers)

	go func() {
		if err := h.startMessageLoop(ctx); err != nil && err != context.Canceled {
			logger.Error("Message loop error: %v", err)
		}
	}()

	return nil
}

func (h *HotkeyManager) startMessageLoop(ctx context.Context) error {
	var msg MSG
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := h.processMessage(&msg); err != nil {
				if e, ok := err.(syscall.Errno); !ok || e != ERROR_SUCCESS {
					logger.Error("Error processing messages: %v", err)
				}
			}
		}
	}
}

func (h *HotkeyManager) processMessage(msg *MSG) error {
	success, err := peekMessage(msg)
	if !success {
		return nil
	}

	if msg.Message == WM_HOTKEY {
		select {
		case h.eventChan <- struct{}{}:
		default:
			logger.Debug("Skipping hotkey event - channel full")
		}
	}

	return err
}

// GetEventChannel returns the channel that emits hotkey events
func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

// Cleanup performs any necessary cleanup
func (h *HotkeyManager) Cleanup() error {
	success, err := unregisterHotkey(h.config.WindowHandle, h.config.ID)
	if !success {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}
	return nil
}
