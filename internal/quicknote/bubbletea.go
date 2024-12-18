package quicknote

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type BubbleTeaQuickNote struct {
	program *tea.Program
	input   chan string
}

func newBubbleTeaQuickNote() (*BubbleTeaQuickNote, error) {
	input := make(chan string)
	program := tea.NewProgram(initialModel())

	return &BubbleTeaQuickNote{
		program: program,
		input:   input,
	}, nil
}

func (b *BubbleTeaQuickNote) Show(ctx context.Context) error {
	// Run the Bubble Tea program
	go func() {
		if err := b.program.Start(); err != nil {
			// Handle error
		}
	}()
	return nil
}

func (b *BubbleTeaQuickNote) GetInput() <-chan string {
	return b.input
}

// Add necessary model implementation for Bubble Tea
type model struct {
	input string
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return "Quick Note: " + m.input
}
