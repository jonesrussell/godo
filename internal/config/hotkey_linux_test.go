//go:build linux && !windows && !darwin && !docker
// +build linux,!windows,!darwin,!docker

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jonesrussell/godo/internal/config"
)

func TestLinuxGetDefaultQuickNoteKey(t *testing.T) {
	key := config.GetDefaultQuickNoteKey()
	assert.Empty(t, key, "Default quick note key should be empty")
}

func TestLinuxGetDefaultQuickNoteKeyString(t *testing.T) {
	keyStr := config.GetDefaultQuickNoteKeyString()
	assert.Empty(t, keyStr, "Default quick note key should be empty")
}

func TestLinuxGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := config.GetDefaultQuickNoteHotkey()
	assert.Empty(t, hotkey.Key, "Default hotkey key should be empty")
	assert.Empty(t, hotkey.Modifiers, "Default hotkey modifiers should be empty")
}

func TestLinuxGetDefaultQuickNoteModifiersString(t *testing.T) {
	modString := config.GetDefaultQuickNoteModifiersString()
	assert.Empty(t, modString, "Default modifiers string should be empty")
}

func TestLinuxGetDefaultQuickNoteCombo(t *testing.T) {
	combo := config.GetDefaultQuickNoteCombo()
	assert.Empty(t, combo, "Default quick note combo should be empty")
}
