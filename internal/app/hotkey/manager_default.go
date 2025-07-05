//go:build !windows && !linux && !darwin
// +build !windows,!linux,!darwin

package hotkey

import (
	"errors"
	
	"github.com/jonesrussell/godo/internal/common"
)

type platformManager struct {
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
	return errors.New("hotkeys are not supported on this platform")
}

func (m *platformManager) Unregister() error {
	return nil
}

func (m *platformManager) Start() error {
	return errors.New("hotkeys are not supported on this platform")
}

func (m *platformManager) Stop() error {
	return nil
}
