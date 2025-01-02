//go:build linux && !windows && !darwin && !docker
// +build linux,!windows,!darwin,!docker

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinuxGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.Equal(t, "G", key, "Default quick note key should be 'G'")
}

func TestLinuxGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := GetDefaultQuickNoteKeyString()
	assert.Equal(t, "G", keyStr, "Default quick note key should be 'G'")
}

func TestLinuxGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := GetDefaultQuickNoteHotkey()
	assert.Equal(t, "G", hotkey.Key, "Default hotkey key should be 'G'")
	assert.Equal(t, []string{"ctrl", "shift"}, hotkey.Modifiers, "Default hotkey modifiers should be ctrl+shift")
}

func TestLinuxGetDefaultQuickNoteModifiersString(t *testing.T) {
	modString := GetDefaultQuickNoteModifiersString()
	assert.Equal(t, "Ctrl+Shift", modString, "Default modifiers string should be 'Ctrl+Shift'")
}

func TestLinuxGetDefaultQuickNoteCombo(t *testing.T) {
	combo := GetDefaultQuickNoteCombo()
	assert.Equal(t, "Ctrl+Shift+G", combo.String(), "Default quick note combo should be 'Ctrl+Shift+G'")
}
