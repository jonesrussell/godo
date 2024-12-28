//go:build linux && !windows && !darwin
// +build linux,!windows,!darwin

package hotkey

import (
	"golang.design/x/hotkey"
)

type platformManager struct {
	hk        *hotkey.Hotkey
	quickNote QuickNoteService
}

func newPlatformManager(quickNote QuickNoteService) Manager {
	return &platformManager{
		quickNote: quickNote,
	}
}

func (m *platformManager) Register() error {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyN)
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
