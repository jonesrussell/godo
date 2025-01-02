package config

import (
	"strings"

	"github.com/jonesrussell/godo/internal/common"
)

// Modifier represents a hotkey modifier (Ctrl, Alt, etc.)
type Modifier uint8

// Key represents a keyboard key
type Key uint16

// Event represents a hotkey event
type Event struct{}

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register() error
	Unregister() error
	Keydown() <-chan Event
}

// HotkeyFactory creates platform-specific hotkey handlers
type HotkeyFactory interface {
	NewHotkey([]Modifier, Key) HotkeyHandler
}

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	QuickNote common.HotkeyBinding `mapstructure:"quick_note"`
}

// HotkeyCombo represents a hotkey combination as a string
type HotkeyCombo string

// NewHotkeyCombo creates a new hotkey string from modifiers and key
func NewHotkeyCombo(modifiers []string, key string) HotkeyCombo {
	return HotkeyCombo(strings.Join(append(modifiers, key), "+"))
}

// String returns the string representation of the hotkey combination
func (h HotkeyCombo) String() string {
	return string(h)
}

// ParseHotkeyCombo parses a string like "Ctrl+Shift+G" into a HotkeyBinding
func ParseHotkeyCombo(combo string) common.HotkeyBinding {
	parts := strings.Split(combo, "+")
	if len(parts) < 2 {
		return common.HotkeyBinding{}
	}
	return common.HotkeyBinding{
		Modifiers: parts[:len(parts)-1],
		Key:       parts[len(parts)-1],
	}
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (h *HotkeyConfig) UnmarshalText(text []byte) error {
	str := string(text)
	h.QuickNote = ParseHotkeyCombo(str)
	return nil
}
