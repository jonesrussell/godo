//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"sync"

	"github.com/jonesrussell/godo/internal/common"
	"golang.design/x/hotkey"
)

// hotkeyWrapper wraps hotkey.Hotkey to implement HotkeyHandler
type hotkeyWrapper struct {
	hk       *hotkey.Hotkey
	mu       sync.RWMutex
	isActive bool
}

// newHotkeyWrapper creates a new hotkeyWrapper instance
func newHotkeyWrapper(mods []hotkey.Modifier, key hotkey.Key) *hotkeyWrapper {
	return &hotkeyWrapper{
		hk:       hotkey.New(mods, key),
		isActive: false,
	}
}

func (h *hotkeyWrapper) Register(binding *common.HotkeyBinding) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.isActive {
		return nil
	}

	if err := h.hk.Register(); err != nil {
		return err
	}

	h.isActive = true
	return nil
}

func (h *hotkeyWrapper) Unregister(binding *common.HotkeyBinding) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isActive {
		return nil
	}

	if err := h.hk.Unregister(); err != nil {
		return err
	}

	h.isActive = false
	return nil
}

func (h *hotkeyWrapper) Start() error {
	return nil
}

func (h *hotkeyWrapper) Stop() error {
	return h.Unregister(nil)
}

func (h *hotkeyWrapper) IsRegistered() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isActive
}
