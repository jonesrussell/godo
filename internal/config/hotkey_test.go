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
