package hotkey

import "errors"

// Common errors
var (
	ErrUnsupportedPlatform = errors.New("unsupported platform")
)

// BaseHotkeyConfig defines the common configuration for hotkeys
type BaseHotkeyConfig struct {
	ID        int
	Modifiers uint
	Key       rune
}
