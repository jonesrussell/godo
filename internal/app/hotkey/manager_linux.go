//go:build linux && !windows && !darwin
// +build linux,!windows,!darwin

package hotkey

import (
	"fmt"
	"strings"

	"github.com/jonesrussell/godo/internal/common"
	"golang.design/x/hotkey"
)

type platformManager struct {
	hk        *hotkey.Hotkey
	quickNote QuickNoteService
	binding   *common.HotkeyBinding
}

func newPlatformManager(quickNote QuickNoteService, binding *common.HotkeyBinding) Manager {
	return &platformManager{
		quickNote: quickNote,
		binding:   binding,
	}
}

func (m *platformManager) Register() error {
	// Convert string modifiers to hotkey.Modifier
	var mods []hotkey.Modifier
	for _, mod := range m.binding.Modifiers {
		switch strings.ToLower(mod) {
		case "ctrl":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt":
			mods = append(mods, hotkey.ModAlt)
		}
	}

	// Convert key string to hotkey.Key
	var key hotkey.Key
	switch strings.ToUpper(m.binding.Key) {
	case "N":
		key = hotkey.KeyN
	// Add more key mappings as needed
	default:
		return fmt.Errorf("unsupported key: %s", m.binding.Key)
	}

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return err
	}
	m.hk = hk

	// Start listening for hotkey in a goroutine
	go func() {
		for range hk.Keydown() {
			if m.quickNote != nil {
				m.quickNote.Show()
			}
		}
	}()

	return nil
}

func (m *platformManager) Unregister() error {
	if m.hk != nil {
		return m.hk.Unregister()
	}
	return nil
}
