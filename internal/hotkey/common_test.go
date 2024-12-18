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
		{"MOD_ALT", MOD_ALT, 0x0001},
		{"MOD_CONTROL", MOD_CONTROL, 0x0002},
		{"MOD_SHIFT", MOD_SHIFT, 0x0004},
		{"MOD_WIN", MOD_WIN, 0x0008},
		{"MOD_NOREPEAT", MOD_NOREPEAT, 0x4000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%04X, want 0x%04X", tt.name, tt.got, tt.expected)
			}
		})
	}
}
