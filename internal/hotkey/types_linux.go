//go:build linux

package hotkey

// DefaultConfig provides default hotkey configuration for Linux
var DefaultConfig = BaseHotkeyConfig{
	ID:        1,
	Modifiers: MOD_CONTROL | MOD_ALT,
	Key:       'G',
}
