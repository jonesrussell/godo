package quicknote

import "context"

// UI defines the interface for platform-specific quick note UI
type UI interface {
	Show(ctx context.Context) error
	GetInput() <-chan string
}

// Variable to hold the platform-specific quick note UI constructor
var platformConstructor = func() (UI, error) {
	return &defaultQuickNoteUI{
		input: make(chan string),
	}, nil
}

// New creates a new platform-specific quick note UI
func New() (UI, error) {
	return platformConstructor()
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
