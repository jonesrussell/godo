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

// hotkeyWrapper wraps the hotkey package's types to match our interfaces
type hotkeyWrapper struct {
	*hotkey.Hotkey
	events chan config.Event
}

func (h *hotkeyWrapper) Keydown() <-chan config.Event {
	return h.events
}

// NewHotkey creates a new hotkey instance
func (f *defaultHotkeyFactory) NewHotkey(mods []config.Modifier, key config.Key) config.HotkeyHandler {
	// Convert our types to hotkey package types
	hotkeyMods := make([]hotkey.Modifier, len(mods))
	for i, mod := range mods {
		// Safe conversion since both types are uint8
		hotkeyMods[i] = hotkey.Modifier(mod)
	}

	// Safe conversion since both types are uint16
	h := hotkey.New(hotkeyMods, hotkey.Key(key))
	wrapper := &hotkeyWrapper{
		Hotkey: h,
		events: make(chan config.Event, 1),
	}

	// Start goroutine to convert hotkey.Event to config.Event
	go func() {
		for range h.Keydown() {
			wrapper.events <- config.Event{}
		}
	}()

	return wrapper
}
