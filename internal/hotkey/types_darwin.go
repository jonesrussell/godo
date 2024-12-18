//go:build darwin

package hotkey

// DefaultConfig provides default hotkey configuration for macOS
var DefaultConfig = BaseHotkeyConfig{
	ID:        1,
	Modifiers: MOD_COMMAND | MOD_OPTION,
	Key:       'G',
}
