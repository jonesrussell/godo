//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"github.com/stretchr/testify/mock"
	"golang.design/x/hotkey"
)

// hotkeyInterface defines the interface for hotkey functionality
type hotkeyInterface interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
}

type mockHotkey struct {
	mock.Mock
	keydownChan chan hotkey.Event
}

func newMockHotkey(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface {
	return &mockHotkey{
		keydownChan: make(chan hotkey.Event),
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
