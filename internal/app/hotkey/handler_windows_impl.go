//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import (
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/common"
	"golang.design/x/hotkey"
)

// windowsHotkeyHandler implements WindowsHotkeyHandler using golang.design/x/hotkey
type windowsHotkeyHandler struct {
	mu         sync.RWMutex
	hk         *hotkey.Hotkey
	registered bool
}

// NewWindowsHotkeyHandler creates a new Windows hotkey handler
func NewWindowsHotkeyHandler() WindowsHotkeyHandler {
	return &windowsHotkeyHandler{}
}

// Register registers the hotkey with the system
func (h *windowsHotkeyHandler) Register(binding *common.HotkeyBinding) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Convert binding to hotkey modifiers and key
	modifiers := []hotkey.Modifier{}
	for _, mod := range binding.Modifiers {
		switch mod {
		case "Ctrl":
			modifiers = append(modifiers, hotkey.ModCtrl)
		case "Alt":
			modifiers = append(modifiers, hotkey.ModAlt)
		case "Shift":
			modifiers = append(modifiers, hotkey.ModShift)
		default:
			return fmt.Errorf("unsupported modifier: %s", mod)
		}
	}

	// Convert key string to hotkey.Key
	var key hotkey.Key
	switch binding.Key {
	case "Enter":
		key = hotkey.KeyReturn
	case "Space":
		key = hotkey.KeySpace
	case "Escape":
		key = hotkey.KeyEscape
	default:
		// Assume single character keys
		if len(binding.Key) != 1 {
			return fmt.Errorf("unsupported key: %s", binding.Key)
		}
		key = hotkey.Key(binding.Key[0])
	}

	// Create and register hotkey
	h.hk = hotkey.New(modifiers, key)
	if err := h.hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	h.registered = true
	return nil
}

// Unregister removes the hotkey registration from the system
func (h *windowsHotkeyHandler) Unregister(binding *common.HotkeyBinding) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.hk == nil {
		return nil
	}

	if err := h.hk.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister hotkey: %w", err)
	}

	h.registered = false
	h.hk = nil
	return nil
}

// Start begins listening for hotkey events
func (h *windowsHotkeyHandler) Start() error {
	return nil
}

// Stop ends the hotkey listening
func (h *windowsHotkeyHandler) Stop() error {
	return nil
}

// IsRegistered returns whether the hotkey is currently registered
func (h *windowsHotkeyHandler) IsRegistered() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.registered
}
