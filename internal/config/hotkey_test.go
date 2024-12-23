package config

import (
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"golang.design/x/hotkey"
)

func TestHotkey(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name          string
		hotkeyStr     Hotkey
		expectedMods  []hotkey.Modifier
		expectedKey   hotkey.Key
		shouldSucceed bool
	}{
		{
			name:      "Ctrl+Alt+G",
			hotkeyStr: "Ctrl+Alt+G",
			expectedMods: []hotkey.Modifier{
				hotkey.ModCtrl,
				hotkey.ModAlt,
			},
			expectedKey:   hotkey.KeyG,
			shouldSucceed: true,
		},
		{
			name:      "Ctrl+Space",
			hotkeyStr: "Ctrl+Space",
			expectedMods: []hotkey.Modifier{
				hotkey.ModCtrl,
			},
			expectedKey:   hotkey.KeySpace,
			shouldSucceed: true,
		},
		{
			name:          "Invalid hotkey",
			hotkeyStr:     "Invalid+X",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String() method
			assert.Equal(t, string(tt.hotkeyStr), tt.hotkeyStr.String())

			// Test parsing
			mods, key, err := tt.hotkeyStr.Parse(mapper)

			if !tt.shouldSucceed {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, len(tt.expectedMods), len(mods))

			// Check if all expected modifiers are present
			for i, mod := range tt.expectedMods {
				assert.Equal(t, mod, mods[i])
			}
		})
	}
}

func TestHotkeyConfig(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name          string
		config        HotkeyConfig
		expectedHK    Hotkey
		shouldSucceed bool
	}{
		{
			name: "Valid quick note hotkey",
			config: HotkeyConfig{
				QuickNote: "Ctrl+Alt+G",
			},
			expectedHK:    "Ctrl+Alt+G",
			shouldSucceed: true,
		},
		{
			name: "Empty quick note hotkey",
			config: HotkeyConfig{
				QuickNote: "",
			},
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, string(tt.config.QuickNote), tt.config.QuickNote.String())

			if tt.shouldSucceed {
				assert.Equal(t, tt.expectedHK, tt.config.QuickNote)
				assert.True(t, tt.config.QuickNote.IsValid(mapper))
			}
		})
	}
}

func TestHotkey_Parse(t *testing.T) {
	log := logger.NewTestLogger(t)
	mapper := NewHotkeyMapper(log)

	tests := []struct {
		name          string
		hotkeyStr     Hotkey
		expectedMods  []hotkey.Modifier
		expectedKey   hotkey.Key
		shouldSucceed bool
	}{
		{
			name:      "Ctrl+Alt+G",
			hotkeyStr: "Ctrl+Alt+G",
			expectedMods: []hotkey.Modifier{
				hotkey.ModCtrl,
				hotkey.ModAlt,
			},
			expectedKey:   hotkey.KeyG,
			shouldSucceed: true,
		},
		{
			name:      "Ctrl+Space",
			hotkeyStr: "Ctrl+Space",
			expectedMods: []hotkey.Modifier{
				hotkey.ModCtrl,
			},
			expectedKey:   hotkey.KeySpace,
			shouldSucceed: true,
		},
		{
			name:      "Ctrl+Shift+Alt+A",
			hotkeyStr: "Ctrl+Shift+Alt+A",
			expectedMods: []hotkey.Modifier{
				hotkey.ModCtrl,
				hotkey.ModShift,
				hotkey.ModAlt,
			},
			expectedKey:   hotkey.KeyA,
			shouldSucceed: true,
		},
		{
			name:          "Invalid hotkey",
			hotkeyStr:     "Invalid+Key+Combo",
			shouldSucceed: false,
		},
		{
			name:          "Empty hotkey",
			hotkeyStr:     "",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String() method
			assert.Equal(t, string(tt.hotkeyStr), tt.hotkeyStr.String())

			// Test parsing
			mods, key, err := tt.hotkeyStr.Parse(mapper)

			if !tt.shouldSucceed {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedKey, key)
			assert.Equal(t, len(tt.expectedMods), len(mods))

			// Check if all expected modifiers are present
			for i, mod := range tt.expectedMods {
				assert.Equal(t, mod, mods[i])
			}

			// Test roundtrip
			newHK := NewHotkey(mapper, mods, key)
			assert.Equal(t, tt.hotkeyStr, newHK)
		})
	}
}
