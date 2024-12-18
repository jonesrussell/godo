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

func (m *Manager) Start(ctx context.Context) error {
	// Cleanup any existing registration
	_, _ = unregisterHotkey(m.config.WindowHandle, m.config.ID)
	time.Sleep(100 * time.Millisecond)

	success, err := registerHotkey(m.config)
	if !success {
		lastErr := syscall.GetLastError()
		return fmt.Errorf("failed to register hotkey: %w (lastErr=%d)", err, lastErr)
	}

	logger.Info("Successfully registered hotkey (ID=%d, Key='%c', Mods=0x%X)",
		m.config.ID, m.config.Key, m.config.Modifiers)

	go func() {
		if err := m.startMessageLoop(ctx); err != nil && err != context.Canceled {
			logger.Error("Message loop error: %v", err)
		}
	}()

	return nil
}

func (m *Manager) startMessageLoop(ctx context.Context) error {
	var msg MSG
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.processMessage(&msg); err != nil {
				if e, ok := err.(syscall.Errno); !ok || e != ERROR_SUCCESS {
					logger.Error("Error processing messages: %v", err)
				}
			}
		}
	}
}

func (m *Manager) processMessage(msg *MSG) error {
	success, err := peekMessage(msg)
	if !success {
		return nil
	}

	if msg.Message == WM_HOTKEY {
		select {
		case m.hotkeyPressed <- struct{}{}:
		default:
			logger.Debug("Skipping hotkey event - channel full")
		}
	}

	return err
}

func (m *Manager) Cleanup() error {
	success, err := unregisterHotkey(m.config.WindowHandle, m.config.ID)
	if !success {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}
	return nil
}
