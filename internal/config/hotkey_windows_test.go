//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowsGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := GetDefaultQuickNoteHotkey()
	assert.NotNil(t, hotkey, "Default quick note hotkey should not be nil")
	assert.Equal(t, "G", hotkey.Key, "Default quick note key should be 'G'")
	assert.Equal(t, []string{"ctrl", "shift"}, hotkey.Modifiers, "Default quick note modifiers should be ctrl+shift")
}

func TestWindowsGetDefaultQuickNoteCombo(t *testing.T) {
	combo := GetDefaultQuickNoteCombo()
	assert.NotEmpty(t, combo, "Default quick note combo should not be empty")
	assert.Equal(t, "Ctrl+Shift+G", combo.String(), "Default quick note combo should be 'Ctrl+Shift+G'")
}

func TestWindowsGetDefaultQuickNoteModifiersString(t *testing.T) {
	modStr := GetDefaultQuickNoteModifiersString()
	assert.NotEmpty(t, modStr, "Default quick note modifiers string should not be empty")
	assert.Equal(t, "Ctrl+Shift", modStr, "Default quick note modifiers should be 'Ctrl+Shift'")
}

func TestWindowsGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.NotEmpty(t, key, "Default quick note key should not be empty")
	assert.Equal(t, "G", key, "Default quick note key should be 'G'")
}

func TestWindowsGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := GetDefaultQuickNoteKeyString()
	assert.NotEmpty(t, keyStr, "Default quick note key string should not be empty")
	assert.Equal(t, "G", keyStr, "Default quick note key should be 'G'")
}
