//go:build darwin

package ui

import (
	"context"
)

// DarwinQuickNote implements QuickNoteUI for macOS
type DarwinQuickNote struct {
	base *BubbleTeaQuickNote
}

// newPlatformQuickNoteUI creates a new macOS-specific quick note UI
func newPlatformQuickNoteUI() (QuickNoteUI, error) {
	base, err := newBubbleTeaQuickNote()
	if err != nil {
		return nil, err
	}
	return &DarwinQuickNote{base: base}, nil
}

// Show displays the quick note UI
func (d *DarwinQuickNote) Show(ctx context.Context) error {
	return d.base.Show(ctx)
}

// GetInput returns the input channel
func (d *DarwinQuickNote) GetInput() <-chan string {
	return d.base.GetInput()
}
