//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"github.com/stretchr/testify/mock"
	"golang.design/x/hotkey"
)

type mockHotkey struct {
	mock.Mock
	keydownChan chan hotkey.Event
	modifiers   []hotkey.Modifier
	key         hotkey.Key
}

func newMockHotkey(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface {
	return &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   mods,
		key:         key,
	}
}

func (m *mockHotkey) Register() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHotkey) Unregister() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHotkey) Keydown() <-chan hotkey.Event {
	return m.keydownChan
}

func (m *mockHotkey) SimulateKeyPress() {
	m.keydownChan <- hotkey.Event{}
}
