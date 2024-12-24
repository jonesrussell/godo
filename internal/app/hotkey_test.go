//go:build !docker
// +build !docker

package app

import (
	"github.com/jonesrussell/godo/internal/config"
)

// TestHotkey is a mock hotkey implementation for testing
type TestHotkey struct {
	keydown chan config.Event
}

// Ensure TestHotkey implements config.HotkeyHandler
var _ config.HotkeyHandler = (*TestHotkey)(nil)

func NewTestHotkey() *TestHotkey {
	return &TestHotkey{
		keydown: make(chan config.Event, 1),
	}
}

func (h *TestHotkey) Register() error {
	return nil
}

func (h *TestHotkey) Unregister() error {
	return nil
}

func (h *TestHotkey) Keydown() <-chan config.Event {
	return h.keydown
}

// Trigger simulates a hotkey press
func (h *TestHotkey) Trigger() {
	h.keydown <- config.Event{}
}
