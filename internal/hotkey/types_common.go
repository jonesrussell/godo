package hotkey

// BaseHotkeyConfig defines the common configuration for a hotkey
type BaseHotkeyConfig struct {
	ID        int
	Modifiers uint
	Key       rune
}
