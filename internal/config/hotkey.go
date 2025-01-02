package config

import "github.com/jonesrussell/godo/internal/common"

// HotkeyString represents a hotkey configuration as strings
type HotkeyString struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers"`
}

// HotkeyDefaultsService defines the interface for platform-specific hotkey defaults
type HotkeyDefaultsService interface {
	GetDefaultQuickNoteHotkey() *common.HotkeyBinding
	GetDefaultQuickNoteCombo() string
	GetDefaultQuickNoteModifiersString() string
	GetDefaultQuickNoteKey() string
	GetDefaultQuickNoteKeyString() string
}
