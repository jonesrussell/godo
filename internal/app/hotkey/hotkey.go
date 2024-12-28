// Package hotkey provides global hotkey functionality for the application
package hotkey

// Manager defines the interface for global hotkey functionality
type Manager interface {
	Register() error
	Unregister() error
}

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}

// New creates a new platform-specific hotkey manager
func New(quickNote QuickNoteService) Manager {
	return newPlatformManager(quickNote)
}
