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
	assert.Empty(t, hotkey.Key, "Default quick note key should be empty")
	assert.Empty(t, hotkey.Modifiers, "Default quick note modifiers should be empty")
}

func TestWindowsGetDefaultQuickNoteCombo(t *testing.T) {
	combo := GetDefaultQuickNoteCombo()
	assert.Empty(t, combo, "Default quick note combo should be empty")
}

func TestWindowsGetDefaultQuickNoteModifiersString(t *testing.T) {
	modStr := GetDefaultQuickNoteModifiersString()
	assert.Empty(t, modStr, "Default quick note modifiers string should be empty")
}

func TestWindowsGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.Empty(t, key, "Default quick note key should be empty")
}

func TestWindowsGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := GetDefaultQuickNoteKeyString()
	assert.Empty(t, keyStr, "Default quick note key string should be empty")
}
