//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"github.com/jonesrussell/godo/internal/common"
	"github.com/stretchr/testify/mock"
	"golang.design/x/hotkey"
)

type mockHotkey struct {
	mock.Mock
	keydownChan chan hotkey.Event
	modifiers   []hotkey.Modifier
	key         hotkey.Key
	registered  bool
}

func newMockHotkey(mods []hotkey.Modifier, key hotkey.Key) HotkeyHandler {
	return &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   mods,
		key:         key,
		registered:  false,
	}
}

func (m *mockHotkey) Register(binding *common.HotkeyBinding) error {
	args := m.Called(binding)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	m.registered = true
	return nil
}

func (m *mockHotkey) Unregister(binding *common.HotkeyBinding) error {
	args := m.Called(binding)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	m.registered = false
	return nil
}

func (m *mockHotkey) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHotkey) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHotkey) Keydown() <-chan hotkey.Event {
	return m.keydownChan
}

// SimulateKeyPress simulates a hotkey press by sending an event
func (m *mockHotkey) SimulateKeyPress() {
	if m.registered {
		m.keydownChan <- hotkey.Event{}
	}
}

// IsRegistered returns whether the hotkey is currently registered
func (m *mockHotkey) IsRegistered() bool {
	return m.registered
}

// Close cleans up resources
func (m *mockHotkey) Close() {
	close(m.keydownChan)
}
