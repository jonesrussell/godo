//go:build !ci && !android && !ios && !wasm && !test_web_driver && !docker

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotkeyString(t *testing.T) {
	tests := []struct {
		name      string
		modifiers []string
		key       string
		want      HotkeyString
	}{
		{
			name:      "linux default",
			modifiers: []string{"Ctrl", "Alt"},
			key:       "G",
			want:      "Ctrl+Alt+G",
		},
		{
			name:      "with shift",
			modifiers: []string{"Ctrl", "Shift", "Alt"},
			key:       "A",
			want:      "Ctrl+Shift+Alt+A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHotkeyString(tt.modifiers, tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDefaultQuickNoteKey(t *testing.T) {
	key := GetDefaultQuickNoteKey()
	assert.NotEmpty(t, key, "Default quick note key should not be empty")
}

func TestGetDefaultQuickNoteModifiers(t *testing.T) {
	modifiers := GetDefaultQuickNoteModifiers()
	assert.NotNil(t, modifiers, "Default quick note modifiers should not be nil")
}

func TestGetDefaultQuickNoteKeyString(t *testing.T) {
	keyString := GetDefaultQuickNoteKeyString()
	assert.NotEmpty(t, keyString, "Default quick note key string should not be empty")
}

func TestGetDefaultQuickNoteModifiersString(t *testing.T) {
	modifiersString := GetDefaultQuickNoteModifiersString()
	assert.NotEmpty(t, modifiersString, "Default quick note modifiers string should not be empty")
}

func TestGetDefaultQuickNoteHotkey(t *testing.T) {
	hotkey := GetDefaultQuickNoteHotkey()
	assert.NotEmpty(t, hotkey, "Default quick note hotkey should not be empty")
}
