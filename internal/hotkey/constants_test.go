package hotkey

import "testing"

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	tests := []struct {
		name     string
		got      uint
		expected uint
	}{
		{"MOD_ALT", MOD_ALT, 0x0001},
		{"MOD_CONTROL", MOD_CONTROL, 0x0002},
		{"WM_HOTKEY", WM_HOTKEY, 0x0312},
		{"PM_REMOVE", PM_REMOVE, 0x0001},
		{"ERROR_HOTKEY_ALREADY_REGISTERED", ERROR_HOTKEY_ALREADY_REGISTERED, 0x0402},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = 0x%04X, want 0x%04X", tt.name, tt.got, tt.expected)
			}
		})
	}
}
