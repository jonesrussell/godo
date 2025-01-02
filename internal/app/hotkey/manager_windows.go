//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
)

// WindowsManager implements the Manager interface for Windows
type WindowsManager struct {
	log       logger.Logger
	quickNote gui.QuickNoteManager
	binding   *common.HotkeyBinding
	hk        HotkeyHandler
	isActive  bool
}

// NewWindowsManager creates a new Windows hotkey manager
func NewWindowsManager(log logger.Logger) (*WindowsManager, error) {
	if log == nil {
		return nil, fmt.Errorf("logger is required")
	}
	return &WindowsManager{
		log: log,
	}, nil
}

// SetQuickNote sets the quick note service and hotkey binding for this manager.
// Both the service and binding are required for the hotkey to function.
func (m *WindowsManager) SetQuickNote(quickNote gui.QuickNoteManager, binding *common.HotkeyBinding) {
	m.log.Info("Setting quick note service and binding",
		"binding", binding,
		"quicknote_nil", quickNote == nil)
	m.quickNote = quickNote
	m.binding = binding
}

// SetHotkey sets the hotkey handler for this manager
func (m *WindowsManager) SetHotkey(hk HotkeyHandler) {
	m.hk = hk
}

// Register registers the hotkey with the system
func (m *WindowsManager) Register() error {
	m.log.Info("Registering hotkey",
		"binding", m.binding,
		"hotkey_nil", m.hk == nil)

	if m.hk == nil {
		m.log.Error("Cannot register - hotkey handler not set")
		return fmt.Errorf("hotkey handler not set")
	}

	if m.binding == nil {
		m.log.Error("Cannot register - binding not set")
		return fmt.Errorf("binding not set")
	}

	if err := m.hk.Register(m.binding); err != nil {
		m.log.Error("Failed to register hotkey",
			"error", err,
			"binding", m.binding)
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	m.isActive = true
	m.log.Info("Hotkey registered successfully",
		"binding", m.binding)
	return nil
}

// Unregister unregisters the hotkey from the system
func (m *WindowsManager) Unregister() error {
	m.log.Info("Unregistering hotkey",
		"binding", m.binding,
		"hotkey_nil", m.hk == nil,
		"is_active", m.isActive)

	if !m.isActive {
		m.log.Warn("Hotkey not active - skipping unregister")
		return nil
	}

	if m.hk == nil {
		m.log.Error("Cannot unregister - hotkey handler not set")
		return fmt.Errorf("hotkey handler not set")
	}

	if m.binding == nil {
		m.log.Error("Cannot unregister - binding not set")
		return fmt.Errorf("binding not set")
	}

	if err := m.hk.Unregister(m.binding); err != nil {
		m.log.Error("Failed to unregister hotkey",
			"error", err,
			"binding", m.binding)
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}

	m.isActive = false
	m.log.Info("Hotkey unregistered successfully")
	return nil
}

// Start begins listening for hotkey events and shows the quick note window when triggered.
// Returns an error if the hotkey is not registered or the quick note service is not set.
func (m *WindowsManager) Start() error {
	m.log.Info("Starting hotkey manager",
		"binding", m.binding,
		"hotkey_nil", m.hk == nil,
		"quicknote_nil", m.quickNote == nil)

	if m.quickNote == nil {
		m.log.Error("Cannot start - quick note service not set")
		return fmt.Errorf("quick note service not set")
	}

	if err := m.hk.Start(); err != nil {
		m.log.Error("Failed to start hotkey handler", "error", err)
		return fmt.Errorf("failed to start hotkey handler: %w", err)
	}

	m.log.Info("Hotkey manager started successfully")
	return nil
}

// Stop stops listening for hotkey events
func (m *WindowsManager) Stop() error {
	m.log.Info("Stopping hotkey manager",
		"binding", m.binding,
		"hotkey_nil", m.hk == nil)

	if m.hk == nil {
		m.log.Error("Cannot stop - hotkey handler not set")
		return fmt.Errorf("hotkey handler not set")
	}

	if err := m.hk.Stop(); err != nil {
		m.log.Error("Failed to stop hotkey handler", "error", err)
		return fmt.Errorf("failed to stop hotkey handler: %w", err)
	}

	m.log.Info("Hotkey manager stopped successfully")
	return nil
}

// IsActive returns whether the hotkey is currently active
func (m *WindowsManager) IsActive() bool {
	return m.isActive
}
