//go:build linux && !windows && !darwin && !docker
// +build linux,!windows,!darwin,!docker

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.Empty(t, key, "Default quick note key should be empty")
}

func TestLinuxGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := GetDefaultQuickNoteKeyString()
	assert.Empty(t, keyStr, "Default quick note key should be empty")
}

func TestLinuxGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := GetDefaultQuickNoteHotkey()
	assert.Empty(t, hotkey.Key, "Default hotkey key should be empty")
	assert.Empty(t, hotkey.Modifiers, "Default hotkey modifiers should be empty")
}

func TestLinuxGetDefaultQuickNoteModifiersString(t *testing.T) {
	modString := GetDefaultQuickNoteModifiersString()
	assert.Empty(t, modString, "Default modifiers string should be empty")
}

func TestLinuxGetDefaultQuickNoteCombo(t *testing.T) {
	combo := GetDefaultQuickNoteCombo()
	assert.Empty(t, combo, "Default quick note combo should be empty")
}
