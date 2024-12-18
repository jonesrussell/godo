//go:build linux

package ui

import (
	"context"
)

// LinuxQuickNote implements QuickNoteUI for Linux
type LinuxQuickNote struct {
	base *BubbleTeaQuickNote
}

// newPlatformQuickNoteUI creates a new Linux-specific quick note UI
func newPlatformQuickNoteUI() (QuickNoteUI, error) {
	base, err := newBubbleTeaQuickNote()
	if err != nil {
		return nil, err
	}
	return &LinuxQuickNote{base: base}, nil
}

// Show displays the quick note UI
func (l *LinuxQuickNote) Show(ctx context.Context) error {
	return l.base.Show(ctx)
}

// GetInput returns the input channel
func (l *LinuxQuickNote) GetInput() <-chan string {
	return l.base.GetInput()
}
