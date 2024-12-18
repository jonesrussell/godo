package ui

import (
	"context"
)

// QuickNoteUI defines the interface for platform-specific quick note implementations
type QuickNoteUI interface {
	Show(ctx context.Context) error
	GetInput() <-chan string
}

// NewQuickNoteUI creates a new platform-specific quick note UI
func NewQuickNoteUI() (QuickNoteUI, error) {
	return newPlatformQuickNoteUI()
}
