// Package hotkey provides hotkey management functionality for the application
package hotkey

import (
	"golang.design/x/hotkey"
)

// Manager defines the interface for hotkey management
type Manager interface {
	// Register registers the hotkey with the system
	Register() error
	// Unregister removes the hotkey registration from the system
	Unregister() error
	// Start begins listening for hotkey events
	Start() error
	// Stop ends the hotkey listening and unregisters the hotkey
	Stop() error
	// IsActive returns whether the hotkey is currently active
	IsActive() bool
}

// DefaultManager implements Manager using golang.design/x/hotkey
type DefaultManager struct {
	hk       *hotkey.Hotkey
	isActive bool
}

// NewManager creates a new DefaultManager instance
func NewManager(modifiers []hotkey.Modifier, key hotkey.Key) (*DefaultManager, error) {
	hk := hotkey.New(modifiers, key)
	return &DefaultManager{
		hk: hk,
	}, nil
}

// Register registers the hotkey with the system
func (m *DefaultManager) Register() error {
	if err := m.hk.Register(); err != nil {
		return err
	}
	m.isActive = true
	return nil
}

// Unregister removes the hotkey registration from the system
func (m *DefaultManager) Unregister() error {
	if err := m.hk.Unregister(); err != nil {
		return err
	}
	m.isActive = false
	return nil
}

// Start begins listening for hotkey events
func (m *DefaultManager) Start() error {
	return nil
}

// Stop ends the hotkey listening and unregisters the hotkey
func (m *DefaultManager) Stop() error {
	return m.Unregister()
}

// IsActive returns whether the hotkey is currently active
func (m *DefaultManager) IsActive() bool {
	return m.isActive
}

// GetHotkey returns the underlying hotkey instance
func (m *DefaultManager) GetHotkey() *hotkey.Hotkey {
	return m.hk
}
