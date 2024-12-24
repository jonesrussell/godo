package config

import (
	"golang.design/x/hotkey"
)

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
}

// HotkeyFactory creates platform-specific hotkey handlers
type HotkeyFactory interface {
	NewHotkey([]hotkey.Modifier, hotkey.Key) HotkeyHandler
}
