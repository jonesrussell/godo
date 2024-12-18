//go:build linux
// +build linux

package quicknote

import "context"

func init() {
	platformConstructor = func() (QuickNoteUI, error) {
		base, err := newBubbleTeaQuickNote()
		if err != nil {
			return nil, err
		}
		return &LinuxQuickNote{base: base}, nil
	}
}

type LinuxQuickNote struct {
	base *BubbleTeaQuickNote
}

func (l *LinuxQuickNote) Show(ctx context.Context) error {
	return l.base.Show(ctx)
}

func (l *LinuxQuickNote) GetInput() <-chan string {
	return l.base.GetInput()
}
