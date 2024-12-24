package config

import (
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"golang.design/x/hotkey"
)

func TestHotkeyParsing(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name    string
		hotkey  HotkeyString
		wantErr bool
	}{
		{
			name:    "valid hotkey",
			hotkey:  "Ctrl+Alt+G",
			wantErr: false,
		},
		{
			name:    "invalid modifier",
			hotkey:  "Invalid+G",
			wantErr: true,
		},
		{
			name:    "invalid key",
			hotkey:  "Ctrl+Invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			hotkey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tt.hotkey.Parse(mapper)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHotkeyCreation(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name     string
		mods     []hotkey.Modifier
		key      hotkey.Key
		expected HotkeyString
	}{
		{
			name:     "ctrl+alt+g",
			mods:     []hotkey.Modifier{ModCtrl, ModAlt},
			key:      hotkey.KeyG,
			expected: "Ctrl+Alt+G",
		},
		{
			name:     "ctrl+shift+a",
			mods:     []hotkey.Modifier{ModCtrl, ModShift},
			key:      hotkey.KeyA,
			expected: "Ctrl+Shift+A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewHotkey(mapper, tt.mods, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHotkeyValidation(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name     string
		hotkey   HotkeyString
		expected bool
	}{
		{
			name:     "valid hotkey",
			hotkey:   "Ctrl+Alt+G",
			expected: true,
		},
		{
			name:     "invalid hotkey",
			hotkey:   "Invalid+Key",
			expected: false,
		},
		{
			name:     "empty string",
			hotkey:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.hotkey.IsValid(mapper))
		})
	}
}
