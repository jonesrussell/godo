//go:build !ci && !android && !ios && !wasm && !test_web_driver && !docker

package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jonesrussell/godo/internal/config"
)

func TestHotkeyCombo(t *testing.T) {
	tests := []struct {
		name      string
		modifiers []string
		key       string
		want      string
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
			got := config.NewHotkeyCombo(tt.modifiers, tt.key)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestHotkeyString(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		modifiers []string
		want      config.HotkeyString
	}{
		{
			name:      "simple hotkey",
			key:       "N",
			modifiers: []string{"ctrl", "shift"},
			want: config.HotkeyString{
				Key:       "N",
				Modifiers: []string{"ctrl", "shift"},
			},
		},
		{
			name:      "with alt",
			key:       "A",
			modifiers: []string{"ctrl", "alt"},
			want: config.HotkeyString{
				Key:       "A",
				Modifiers: []string{"ctrl", "alt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.HotkeyString{
				Key:       tt.key,
				Modifiers: tt.modifiers,
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
