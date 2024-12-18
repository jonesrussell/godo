package hotkey

import (
	"testing"
)

func TestModifierConstants(t *testing.T) {
	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{
			name:     "MOD_ALT",
			got:      MOD_ALT,
			expected: 0x0001,
		},
		{
			name:     "MOD_CONTROL",
			got:      MOD_CONTROL,
			expected: 0x0002,
		},
		{
			name:     "MOD_SHIFT",
			got:      MOD_SHIFT,
			expected: 0x0004,
		},
		{
			name:     "MOD_WIN",
			got:      MOD_WIN,
			expected: 0x0008,
		},
		{
			name:     "MOD_NOREPEAT",
			got:      MOD_NOREPEAT,
			expected: 0x4000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%04X, want 0x%04X", tt.name, tt.got, tt.expected)
			}
		})
	}
}
