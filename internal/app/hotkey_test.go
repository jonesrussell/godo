package app

import (
	"golang.design/x/hotkey"
)

// TestHotkey is a mock hotkey implementation for testing
type TestHotkey struct {
	keydown chan hotkey.Event
}

func NewTestHotkey() *TestHotkey {
	return &TestHotkey{
		keydown: make(chan hotkey.Event),
	}
}

func (h *TestHotkey) Register() error {
	return nil
}

func (h *TestHotkey) Unregister() error {
	return nil
}

func (h *TestHotkey) Keydown() <-chan hotkey.Event {
	return h.keydown
}

// Trigger simulates a hotkey press
func (h *TestHotkey) Trigger() {
	h.keydown <- hotkey.Event{}
}
