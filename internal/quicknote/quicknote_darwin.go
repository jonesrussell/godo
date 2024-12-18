//go:build darwin
// +build darwin

package ui

import "context"

func init() {
	newPlatformQuickNoteUI = func() (QuickNoteUI, error) {
		base, err := newBubbleTeaQuickNote()
		if err != nil {
			return nil, err
		}
		return &DarwinQuickNote{base: base}, nil
	}
}

type DarwinQuickNote struct {
	base *BubbleTeaQuickNote
}

func (d *DarwinQuickNote) Show(ctx context.Context) error {
	return d.base.Show(ctx)
}

func (d *DarwinQuickNote) GetInput() <-chan string {
	return d.base.GetInput()
}
