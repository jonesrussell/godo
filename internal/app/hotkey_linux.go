//go:build linux

package app

import (
	"golang.design/x/hotkey"
)

type linuxHotkeyManager struct {
	hk        *hotkey.Hotkey
	quickNote QuickNoteService
}

func NewHotkeyManager(quickNote QuickNoteService) HotkeyManager {
	return &linuxHotkeyManager{
		quickNote: quickNote,
	}
}

func (m *linuxHotkeyManager) Register() error {
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

func (m *linuxHotkeyManager) Unregister() error {
	if m.hk != nil {
		return m.hk.Unregister()
	}
	return nil
}

// Implementation details...
