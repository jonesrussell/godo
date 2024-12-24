package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"golang.design/x/hotkey"
)

// defaultHotkeyFactory is the default implementation of config.HotkeyFactory
type defaultHotkeyFactory struct{}

// NewHotkeyFactory creates a new default hotkey factory
func NewHotkeyFactory() config.HotkeyFactory {
	return &defaultHotkeyFactory{}
}

// NewHotkey creates a new hotkey instance
func (f *defaultHotkeyFactory) NewHotkey(mods []hotkey.Modifier, key hotkey.Key) config.HotkeyHandler {
	return hotkey.New(mods, key)
}
