package hotkey

import (
	"testing"
)

func TestRegisterHotkey(t *testing.T) {
	config := HotkeyConfig{
		WindowHandle: 0,
		ID:           100, // Use different ID from default
		Modifiers:    MOD_ALT,
		Key:          'T',
	}

	success, err := registerHotkey(config)
	if !success {
		t.Errorf("registerHotkey() failed: %v", err)
	}

	// Cleanup
	_, _ = unregisterHotkey(config.WindowHandle, config.ID)
}

func TestUnregisterHotkey(t *testing.T) {
	// First register a hotkey
	config := HotkeyConfig{
		WindowHandle: 0,
		ID:           101, // Use different ID
		Modifiers:    MOD_CONTROL,
		Key:          'Y',
	}

	success, _ := registerHotkey(config)
	if !success {
		t.Skip("Could not register hotkey for unregister test")
	}

	// Test unregistering
	success, err := unregisterHotkey(config.WindowHandle, config.ID)
	if !success {
		t.Errorf("unregisterHotkey() failed: %v", err)
	}
}

func TestPeekMessage(t *testing.T) {
	var msg MSG
	success, err := peekMessage(&msg)

	// We don't care about the result, just that it doesn't crash
	if err != nil && success {
		t.Errorf("peekMessage() unexpected state: success=%v, err=%v", success, err)
	}
}
