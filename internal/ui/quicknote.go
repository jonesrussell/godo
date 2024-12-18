package ui

import "context"

// QuickNoteUI defines the interface for platform-specific quick note implementations
type QuickNoteUI interface {
	// Show displays the quick note input window
	Show(ctx context.Context) error
	// Hide closes the quick note window
	Hide() error
	// GetInput returns the channel that receives user input
	GetInput() <-chan string
}
