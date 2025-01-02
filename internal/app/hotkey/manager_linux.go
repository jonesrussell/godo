//go:build linux
// +build linux

package hotkey

import (
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/common"
)

type linuxManager struct {
	quickNote QuickNoteService
	binding   *common.HotkeyBinding
	running   bool
	mu        sync.Mutex
}

func newPlatformManager(quickNote QuickNoteService, binding *common.HotkeyBinding) Manager {
	return &linuxManager{
		quickNote: quickNote,
		binding:   binding,
	}
}

func (m *linuxManager) Register() error {
	// TODO: Implement Linux-specific hotkey registration using X11/Wayland
	// For now, return a not implemented error
	return fmt.Errorf("hotkey registration not yet implemented for Linux")
}

func (m *linuxManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return nil
	}

	// TODO: Implement Linux-specific hotkey listening
	// For now, return a not implemented error
	return fmt.Errorf("hotkey listening not yet implemented for Linux")
}

func (m *linuxManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	m.running = false
	return nil
}

func (m *linuxManager) Unregister() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	// TODO: Implement Linux-specific hotkey unregistration
	// For now, just mark as not running
	m.running = false
	return nil
}
