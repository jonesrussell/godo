package config

import (
	"fmt"
	"strings"

	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/shared/common"
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

// HotkeyConfig holds hotkey-related configuration
type HotkeyConfig struct {
	QuickNote common.HotkeyBinding `mapstructure:"quick_note"`
	log       logger.Logger
}

// Validate checks if the hotkey configuration is valid
func (c *HotkeyConfig) Validate() error {
	c.log.Debug("Validating hotkey configuration",
		"modifiers", c.QuickNote.Modifiers,
		"key", c.QuickNote.Key)

	if len(c.QuickNote.Modifiers) == 0 {
		return fmt.Errorf("at least one modifier is required")
	}

	// Validate modifiers
	validModifiers := map[string]bool{
		"ctrl":  true,
		"shift": true,
		"alt":   true,
	}

	for _, mod := range c.QuickNote.Modifiers {
		mod = strings.ToLower(mod)
		if !validModifiers[mod] {
			c.log.Error("Invalid modifier",
				"modifier", mod,
				"valid_modifiers", []string{"ctrl", "shift", "alt"})
			return fmt.Errorf("invalid modifier: %s", mod)
		}
	}

	// Validate key
	if c.QuickNote.Key == "" {
		return fmt.Errorf("key is required")
	}

	// Convert key to uppercase for consistency
	c.QuickNote.Key = strings.ToUpper(c.QuickNote.Key)

	c.log.Info("Hotkey configuration validated successfully",
		"modifiers", c.QuickNote.Modifiers,
		"key", c.QuickNote.Key)

	return nil
}

// String returns a string representation of the hotkey
func (c *HotkeyConfig) String() string {
	return fmt.Sprintf("%s+%s", strings.Join(c.QuickNote.Modifiers, "+"), c.QuickNote.Key)
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

const (
	// MinHotkeyParts is the minimum number of parts required in a hotkey combo
	MinHotkeyParts = 2
)

// ParseHotkeyCombo parses a hotkey combo string into a HotkeyBinding
func ParseHotkeyCombo(combo string) common.HotkeyBinding {
	parts := strings.Split(combo, "+")
	if len(parts) < MinHotkeyParts {
		return common.HotkeyBinding{}
	}
	return common.HotkeyBinding{
		Modifiers: parts[:len(parts)-1],
		Key:       parts[len(parts)-1],
	}
}

// UnmarshalText implements encoding.TextUnmarshaler
func (c *HotkeyConfig) UnmarshalText(text []byte) error {
	str := string(text)
	c.QuickNote = ParseHotkeyCombo(str)
	return nil
}
