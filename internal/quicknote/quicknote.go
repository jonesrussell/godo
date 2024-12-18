package quicknote

import "context"

// QuickNoteUI defines the interface for platform-specific quick note UI
type QuickNoteUI interface {
	Show(ctx context.Context) error
	GetInput() <-chan string
}

// Variable to hold the platform-specific quick note UI constructor
var newPlatformQuickNoteUI = func() (QuickNoteUI, error) {
	return &defaultQuickNoteUI{
		input: make(chan string),
	}, nil
}

// NewQuickNoteUI creates a new platform-specific quick note UI
func NewQuickNoteUI() (QuickNoteUI, error) {
	return newPlatformQuickNoteUI()
}

// defaultQuickNoteUI provides a default implementation
type defaultQuickNoteUI struct {
	input chan string
}

func (u *defaultQuickNoteUI) Show(ctx context.Context) error {
	return nil
}

func (u *defaultQuickNoteUI) GetInput() <-chan string {
	return u.input
}
