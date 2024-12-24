package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"golang.design/x/hotkey"
)

// HotkeyManager handles global hotkey registration and events
type HotkeyManager interface {
	Setup() error
}

// NoopHotkeyManager is a no-op implementation for environments that don't support hotkeys
type NoopHotkeyManager struct {
	app *App
}

func NewNoopHotkeyManager(app *App) HotkeyManager {
	return &NoopHotkeyManager{app: app}
}

func (m *NoopHotkeyManager) Setup() error {
	return nil
}

// DefaultHotkeyManager is the default implementation for environments that support hotkeys
type DefaultHotkeyManager struct {
	app *App
}

func NewDefaultHotkeyManager(app *App) HotkeyManager {
	return &DefaultHotkeyManager{app: app}
}

func (m *DefaultHotkeyManager) Setup() error {
	modifiers := config.GetDefaultQuickNoteModifiers()
	key := hotkey.KeyN

	hk := hotkey.New(modifiers, key)
	if err := hk.Register(); err != nil {
		return err
	}

	go func() {
		for range hk.Keydown() {
			m.app.quickNote.Show()
		}
	}()

	return nil
}

// NewHotkeyFactory creates a new default hotkey factory
func NewHotkeyFactory() config.HotkeyFactory {
	return &defaultHotkeyFactory{}
}

// defaultHotkeyFactory is the default implementation of config.HotkeyFactory
type defaultHotkeyFactory struct{}

// NewHotkey creates a new hotkey instance
func (f *defaultHotkeyFactory) NewHotkey(_ []config.Modifier, _ config.Key) config.HotkeyHandler {
	return &noopHotkey{}
}

// noopHotkey is a no-op implementation of config.HotkeyHandler
type noopHotkey struct{}

func (h *noopHotkey) Register() error {
	return nil
}

func (h *noopHotkey) Unregister() error {
	return nil
}

func (h *noopHotkey) Keydown() <-chan config.Event {
	return make(chan config.Event)
}
