package config

// HotkeyString represents a hotkey configuration as strings
type HotkeyString struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers"`
}

// HotkeyDefaults defines the interface for platform-specific hotkey defaults
type HotkeyDefaults interface {
	GetDefaultQuickNoteKey() string
	GetDefaultQuickNoteKeyString() string
	GetDefaultQuickNoteModifiersString() string
	GetDefaultQuickNoteHotkey() HotkeyString
}
