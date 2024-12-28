//go:build !windows && !linux && !darwin
// +build !windows,!linux,!darwin

package hotkey

import "errors"

type platformManager struct {
	quickNote QuickNoteService
}

func newPlatformManager(quickNote QuickNoteService) Manager {
	return &platformManager{
		quickNote: quickNote,
	}
}

func (m *platformManager) Register() error {
	return errors.New("hotkeys are not supported on this platform")
}

func (m *platformManager) Unregister() error {
	return nil
}
