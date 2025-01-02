//go:build linux && !windows && !darwin && !docker
// +build linux,!windows,!darwin,!docker

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.Equal(t, "N", key, "Default quick note key should be 'N'")
}

func TestLinuxGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := GetDefaultQuickNoteKeyString()
	assert.Equal(t, "N", keyStr, "Default quick note key should be 'N'")
}

func TestLinuxGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := GetDefaultQuickNoteHotkey()
	assert.Equal(t, "N", hotkey.Key, "Default hotkey key should be 'N'")
	assert.Equal(t, []string{"ctrl", "shift"}, hotkey.Modifiers, "Default hotkey modifiers should be ctrl+shift")
}

func TestLinuxGetDefaultQuickNoteModifiersString(t *testing.T) {
	modString := GetDefaultQuickNoteModifiersString()
	assert.Equal(t, "Ctrl+Shift", modString, "Default modifiers string should be 'Ctrl+Shift'")
}

func TestLinuxGetDefaultQuickNoteCombo(t *testing.T) {
	combo := GetDefaultQuickNoteCombo()
	assert.Equal(t, "Ctrl+Shift+N", combo.String(), "Default quick note combo should be 'Ctrl+Shift+N'")
}
