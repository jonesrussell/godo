package quicknote

import (
	"context"
)

// UI defines the interface for quick note functionality
type UI interface {
	Show(context.Context) error
	GetInput() <-chan string
	Close()
}

type QuickNote struct {
	input chan string
}

func New() (UI, error) {
	return &QuickNote{
		input: make(chan string, 1),
	}, nil
}

func (qn *QuickNote) Show(ctx context.Context) error {
	// Implementation will be added later with Fyne
	return nil
}

func (qn *QuickNote) GetInput() <-chan string {
	return qn.input
}

func (qn *QuickNote) Close() {
	close(qn.input)
}
