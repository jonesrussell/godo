package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.design/x/hotkey"
)

func TestHotkeyConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    HotkeyConfig
		wantKey   hotkey.Key
		wantMods  []hotkey.Modifier
		wantError bool
	}{
		{
			name: "valid ctrl+space hotkey",
			config: HotkeyConfig{
				Modifiers: []string{"ctrl"},
				Key:       "space",
			},
			wantKey:  hotkey.KeySpace,
			wantMods: []hotkey.Modifier{hotkey.ModCtrl},
		},
		{
			name: "valid ctrl+alt+s hotkey",
			config: HotkeyConfig{
				Modifiers: []string{"ctrl", "alt"},
				Key:       "s",
			},
			wantKey:  hotkey.KeyS,
			wantMods: []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt},
		},
		{
			name: "invalid modifier",
			config: HotkeyConfig{
				Modifiers: []string{"invalid"},
				Key:       "space",
			},
			wantError: true,
		},
		{
			name: "invalid key",
			config: HotkeyConfig{
				Modifiers: []string{"ctrl"},
				Key:       "invalid",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test validation first
			err := tt.config.Validate()

			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Only create hotkey if validation passed
			hk, err := tt.config.ToHotkey()
			require.NoError(t, err)
			require.NotNil(t, hk)

			// Test string representation
			assert.Contains(t, tt.config.String(), tt.config.Key)
			for _, mod := range tt.config.Modifiers {
				assert.Contains(t, tt.config.String(), mod)
			}
		})
	}
}
