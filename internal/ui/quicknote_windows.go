//go:build windows

package ui

import (
	"context"
)

// WindowsQuickNote implements QuickNoteUI for Windows
type WindowsQuickNote struct {
	base *BubbleTeaQuickNote
}

// newPlatformQuickNoteUI creates a new Windows-specific quick note UI
func newPlatformQuickNoteUI() (QuickNoteUI, error) {
	base, err := newBubbleTeaQuickNote()
	if err != nil {
		return nil, err
	}
	return &WindowsQuickNote{base: base}, nil
}

// Show displays the quick note UI
func (w *WindowsQuickNote) Show(ctx context.Context) error {
	return w.base.Show(ctx)
}

// GetInput returns the input channel
func (w *WindowsQuickNote) GetInput() <-chan string {
	return w.base.GetInput()
}
