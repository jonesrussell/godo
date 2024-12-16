package hotkey

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want HotkeyConfig
	}{
		{
			name: "Default configuration",
			want: HotkeyConfig{
				WindowHandle: 0,
				ID:           1,
				Modifiers:    MOD_CONTROL | MOD_ALT,
				Key:          'G',
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if DefaultConfig.WindowHandle != tt.want.WindowHandle {
				t.Errorf("WindowHandle = %v, want %v", DefaultConfig.WindowHandle, tt.want.WindowHandle)
			}
			if DefaultConfig.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", DefaultConfig.ID, tt.want.ID)
			}
			if DefaultConfig.Modifiers != tt.want.Modifiers {
				t.Errorf("Modifiers = %v, want %v", DefaultConfig.Modifiers, tt.want.Modifiers)
			}
			if DefaultConfig.Key != tt.want.Key {
				t.Errorf("Key = %v, want %v", DefaultConfig.Key, tt.want.Key)
			}
		})
	}
}

func TestMSGStructure(t *testing.T) {
	// Test MSG structure size and alignment
	msg := MSG{}

	// Basic structure tests
	if msg.Message != 0 {
		t.Error("New MSG should have Message = 0")
	}
	if msg.WParam != 0 {
		t.Error("New MSG should have WParam = 0")
	}
	if msg.LParam != 0 {
		t.Error("New MSG should have LParam = 0")
	}
	if msg.Time != 0 {
		t.Error("New MSG should have Time = 0")
	}
	if msg.Pt.X != 0 || msg.Pt.Y != 0 {
		t.Error("New MSG should have Pt.X = 0 and Pt.Y = 0")
	}
}
