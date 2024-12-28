// Package hotkey provides global hotkey functionality for the application
package hotkey

import "github.com/jonesrussell/godo/internal/common"

// QuickNoteService defines quick note operations that can be triggered by hotkeys
type QuickNoteService interface {
	Show()
	Hide()
}

// New creates a new platform-specific hotkey manager
func New(quickNote QuickNoteService, binding *common.HotkeyBinding) Manager {
	return newPlatformManager(quickNote, binding)
}
