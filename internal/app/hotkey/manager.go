package hotkey

import (
	"golang.design/x/hotkey"
)

// Manager defines the interface for hotkey management
type Manager interface {
	Register() error
	Unregister() error
	Start() error
	Stop() error
}

// DefaultManager implements Manager using golang.design/x/hotkey
type DefaultManager struct {
	hk *hotkey.Hotkey
}

// NewManager creates a new DefaultManager
func NewManager(modifiers []hotkey.Modifier, key hotkey.Key) (*DefaultManager, error) {
	hk := hotkey.New(modifiers, key)
	if err := hk.Register(); err != nil {
		return nil, err
	}

	return &DefaultManager{
		hk: hk,
	}, nil
}

func (m *DefaultManager) Register() error {
	return m.hk.Register()
}

func (m *DefaultManager) Unregister() error {
	return m.hk.Unregister()
}

func (m *DefaultManager) Start() error {
	go func() {
		for range m.hk.Keydown() {
			// Handle keydown event
		}
	}()
	return nil
}

func (m *DefaultManager) Stop() error {
	return m.hk.Unregister()
}
