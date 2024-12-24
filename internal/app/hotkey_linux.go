//go:build linux

package app

import "golang.design/x/hotkey"

type linuxHotkeyManager struct {
	hk *hotkey.Hotkey
}

func NewHotkeyManager() HotkeyManager {
	return &linuxHotkeyManager{}
}

func (m *linuxHotkeyManager) Register() error {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyN)
	if err := hk.Register(); err != nil {
		return err
	}
	m.hk = hk
	return nil
}

func (m *linuxHotkeyManager) Unregister() error {
	if m.hk != nil {
		return m.hk.Unregister()
	}
	return nil
}

// Implementation details...
