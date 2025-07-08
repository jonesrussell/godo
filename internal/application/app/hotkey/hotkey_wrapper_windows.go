//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"sync"

	"golang.design/x/hotkey"
)

// hotkeyWrapper wraps hotkey.Hotkey to implement hotkeyInterface
type hotkeyWrapper struct {
	hk       *hotkey.Hotkey
	mu       sync.RWMutex
	isActive bool
}

// newHotkeyWrapper creates a new hotkeyWrapper instance
func newHotkeyWrapper(mods []hotkey.Modifier, key hotkey.Key) hotkeyInterface {
	return &hotkeyWrapper{
		hk:       hotkey.New(mods, key),
		isActive: false,
	}
}

func (h *hotkeyWrapper) Register() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := h.hk.Register()
	if err == nil {
		h.isActive = true
	}
	return err
}

func (h *hotkeyWrapper) Unregister() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	err := h.hk.Unregister()
	if err == nil {
		h.isActive = false
	}
	return err
}

func (h *hotkeyWrapper) Keydown() <-chan hotkey.Event {
	return h.hk.Keydown()
}

func (h *hotkeyWrapper) IsRegistered() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isActive
}
