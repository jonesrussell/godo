package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/godo/internal/logger"
)

var (
	style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#383838")).
		PaddingLeft(2)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#383838")).
			PaddingLeft(2).
			PaddingRight(2)

	errorStyle = style.Foreground(lipgloss.Color("#FF0000"))
	hintStyle  = style.Foreground(lipgloss.Color("#666666"))
)

// BubbleTeaQuickNote is the Bubble Tea implementation of QuickNoteUI
type BubbleTeaQuickNote struct {
	input   string
	inputCh chan string
	program *tea.Program
	focused bool
	err     error
}

// Ensure BubbleTeaQuickNote implements QuickNoteUI
var _ QuickNoteUI = (*BubbleTeaQuickNote)(nil)

// newBubbleTeaQuickNote creates a new BubbleTeaQuickNote instance
func newBubbleTeaQuickNote() (*BubbleTeaQuickNote, error) {
	inputCh := make(chan string)
	ui := &BubbleTeaQuickNote{
		inputCh: inputCh,
	}

	program := tea.NewProgram(ui)
	ui.program = program

	return ui, nil
}

// Show displays the quick note UI
func (m *BubbleTeaQuickNote) Show(ctx context.Context) error {
	m.focused = true
	defer func() { m.focused = false }()

	// Create a channel for program completion
	done := make(chan error, 1)

	go func() {
		model, err := m.program.Run()
		if err != nil {
			logger.Error("Error running quick note UI: %v", err)
			m.err = fmt.Errorf("UI error: %w", err)
		}
		// Store final model state if needed
		if model != nil {
			if mn, ok := model.(*BubbleTeaQuickNote); ok {
				m.input = mn.input
			}
		}
		done <- err
	}()

	// Wait for either context cancellation or program completion
	select {
	case <-ctx.Done():
		// Context was cancelled, cleanup
		m.program.Quit()
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// GetInput returns the input channel
func (m *BubbleTeaQuickNote) GetInput() <-chan string {
	return m.inputCh
}

// Tea model implementation methods
func (m *BubbleTeaQuickNote) Init() tea.Cmd {
	return nil
}

func (m *BubbleTeaQuickNote) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.input = ""
			return m, tea.Quit
		case tea.KeyEnter:
			if m.input != "" {
				m.inputCh <- m.input
				m.input = ""
			}
			return m, tea.Quit
		case tea.KeyBackspace:
			if m.input != "" {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if keyMsg.Type == tea.KeyRunes {
				m.input += string(keyMsg.Runes)
			}
		}
	}

	return m, nil
}

func (m *BubbleTeaQuickNote) View() string {
	if !m.focused {
		return ""
	}

	var s string
	s += style.Render("Quick Note") + "\n\n"
	s += inputStyle.Render(m.input)
	if m.err != nil {
		s += "\n" + errorStyle.Render(m.err.Error())
	}
	s += "\n\n" + hintStyle.Render("Press Enter to save, Esc to cancel")
	return s
}
