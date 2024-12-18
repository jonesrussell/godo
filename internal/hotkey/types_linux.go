//go:build linux

package hotkey

// HotkeyConfig defines the configuration for a Linux hotkey
type HotkeyConfig struct {
	ID        int
	Modifiers uint
	Key       rune
}

// DefaultConfig provides a default hotkey configuration
var DefaultConfig = HotkeyConfig{
	ID:        1,
	Modifiers: MOD_CONTROL | MOD_ALT,
	Key:       'G',
}
