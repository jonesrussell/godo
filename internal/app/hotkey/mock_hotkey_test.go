//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"

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

func newMockHotkey(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface {
	return &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   mods,
		key:         key,
		registered:  false,
	}
}

func (m *mockHotkey) Register() error {
	if m.registered {
		return fmt.Errorf("already registered")
	}
	args := m.Called()
	if args.Error(0) == nil {
		m.registered = true
	}
	return args.Error(0)
}

func (m *mockHotkey) Unregister() error {
	if !m.registered {
		return fmt.Errorf("not registered")
	}
	args := m.Called()
	if args.Error(0) == nil {
		m.registered = false
	}
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
