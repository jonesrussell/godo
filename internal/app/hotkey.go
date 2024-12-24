package app

import (
	"golang.design/x/hotkey"
)

// hotkeyInterface defines the interface for hotkey functionality
type hotkeyInterface interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
}

// HotkeyFactory is responsible for creating hotkey instances
type HotkeyFactory interface {
	NewHotkey(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface
}

// defaultHotkeyFactory is the default implementation of HotkeyFactory
type defaultHotkeyFactory struct{}

// NewHotkeyFactory creates a new default hotkey factory
func NewHotkeyFactory() HotkeyFactory {
	return &defaultHotkeyFactory{}
}

// NewHotkey creates a new hotkey instance
func (f *defaultHotkeyFactory) NewHotkey(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface {
	return hotkey.New(mods, key)
}
