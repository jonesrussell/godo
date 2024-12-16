package hotkey

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
)

const (
	ERROR_SUCCESS = 0
)

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	hotkeyPressed chan struct{}
	config        HotkeyConfig
}

// NewHotkeyManager creates a new hotkey manager with default configuration
func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{
		hotkeyPressed: make(chan struct{}),
		config:        DefaultConfig,
	}
}

// GetEventChannel returns the channel for hotkey events
func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.hotkeyPressed
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	logger.Info("Starting hotkey manager...")

	// Cleanup any existing registration
	_, _ = unregisterHotkey(h.config.WindowHandle, h.config.ID)
	time.Sleep(100 * time.Millisecond)

	success, err := registerHotkey(h.config)
	if !success {
		lastErr := syscall.GetLastError()
		return fmt.Errorf("failed to register hotkey: %w (lastErr=%d)", err, lastErr)
	}

	logger.Info("🎮 Successfully registered hotkey (ID=%d, Key='%c', Mods=0x%X)",
		h.config.ID, h.config.Key, h.config.Modifiers)

	// Start message loop in a goroutine
	go func() {
		if err := h.startMessageLoop(ctx); err != nil && err != context.Canceled {
			logger.Error("Message loop error: %v", err)
		}
	}()

	return nil
}

func (h *HotkeyManager) startMessageLoop(ctx context.Context) error {
	logger.Info("🔄 Starting hotkey message loop...")
	var msg MSG
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping hotkey message loop...")
			return ctx.Err()
		case <-ticker.C:
			if err := h.processMessage(&msg); err != nil {
				// Check if it's a Windows success code
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
		return nil // No message available
	}

	if msg.Message == WM_HOTKEY {
		logger.Debug("🎯 Hotkey triggered! (ID=%d)", msg.WParam)
		select {
		case h.hotkeyPressed <- struct{}{}:
			// Message sent successfully
		default:
			// Channel is full, skip this event
			logger.Debug("Skipping hotkey event - channel full")
		}
	}

	return err
}

// Cleanup unregisters the hotkey
func (h *HotkeyManager) Cleanup() error {
	logger.Debug("Cleaning up hotkey registration...")
	success, err := unregisterHotkey(h.config.WindowHandle, h.config.ID)
	if !success {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}
	return nil
}
