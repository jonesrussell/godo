//go:build !windows
// +build !windows

package hotkey

import (
	"context"

	"github.com/jonesrussell/godo/internal/logger"
)

// Manager handles global hotkeys
type Manager struct {
	hotkeyPressed chan struct{}
	config        BaseHotkeyConfig
}

// NewManager creates a new hotkey manager
func NewManager() *Manager {
	return &Manager{
		hotkeyPressed: make(chan struct{}),
		config: BaseHotkeyConfig{
			ID:        1,
			Modifiers: MOD_CONTROL | MOD_ALT,
			Key:       'G',
		},
	}
}

// GetEventChannel returns the channel for hotkey events
func (m *Manager) GetEventChannel() <-chan struct{} {
	return m.hotkeyPressed
}

// Start begins listening for hotkey events
func (m *Manager) Start(ctx context.Context) error {
	logger.Info("Hotkey functionality not supported on this platform")
	return nil
}

// Cleanup performs any necessary cleanup
func (m *Manager) Cleanup() error {
	logger.Debug("Cleaning up hotkey manager (no-op on this platform)")
	return nil
}
