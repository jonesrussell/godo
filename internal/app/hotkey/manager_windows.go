//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"sync"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
)

// WindowsHotkeyManager manages hotkey registrations for Windows
type WindowsHotkeyManager struct {
	mu      sync.RWMutex
	log     logger.Logger
	handler WindowsHotkeyHandler
	binding *common.HotkeyBinding
}

// NewWindowsHotkeyManager creates a new Windows hotkey manager
func NewWindowsHotkeyManager(log logger.Logger, handler WindowsHotkeyHandler) *WindowsHotkeyManager {
	return &WindowsHotkeyManager{
		log:     log,
		handler: handler,
	}
}

// Register registers a hotkey
func (m *WindowsHotkeyManager) Register(binding *common.HotkeyBinding) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.handler.Register(binding); err != nil {
		m.log.Error("failed to register hotkey", "error", err)
		return err
	}

	m.binding = binding
	return nil
}

// Unregister removes the hotkey registration
func (m *WindowsHotkeyManager) Unregister() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.binding == nil {
		return nil
	}

	if err := m.handler.Unregister(m.binding); err != nil {
		m.log.Error("failed to unregister hotkey", "error", err)
		return err
	}

	m.binding = nil
	return nil
}

// Start starts listening for hotkey events
func (m *WindowsHotkeyManager) Start() error {
	return m.handler.Start()
}

// Stop stops listening for hotkey events
func (m *WindowsHotkeyManager) Stop() error {
	return m.handler.Stop()
}

// IsRegistered returns whether a hotkey is currently registered
func (m *WindowsHotkeyManager) IsRegistered() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.binding != nil && m.handler.IsRegistered()
}
