//go:build windows
// +build windows

package ui

import "context"

func init() {
	newPlatformQuickNoteUI = func() (QuickNoteUI, error) {
		base, err := newBubbleTeaQuickNote()
		if err != nil {
			return nil, err
		}
		return &WindowsQuickNote{base: base}, nil
	}
}

type WindowsQuickNote struct {
	base *BubbleTeaQuickNote
}

func (w *WindowsQuickNote) Show(ctx context.Context) error {
	return w.base.Show(ctx)
}

func (w *WindowsQuickNote) GetInput() <-chan string {
	return w.base.GetInput()
}
