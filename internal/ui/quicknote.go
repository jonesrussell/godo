package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
)

type QuickNoteUI struct {
	input   textinput.Model
	service service.TodoServicer
	err     error
}

func NewQuickNote(service service.TodoServicer) *QuickNoteUI {
	input := textinput.New()
	input.Placeholder = "Type your note and press Enter..."
	input.Focus()
	input.Width = 80 // Wider for better visibility

	// Add styling for better visibility
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	return &QuickNoteUI{
		input:   input,
		service: service,
	}
}

func (qn *QuickNoteUI) Init() tea.Cmd {
	logger.Debug("Initializing QuickNote UI")
	return textinput.Blink
}

func (qn *QuickNoteUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		logger.Debug("Received key event: %v (type: %v)", msg.String(), msg.Type)
		switch msg.Type {
		case tea.KeyEnter:
			logger.Debug("Enter key pressed")
			if title := qn.input.Value(); title != "" {
				logger.Debug("Creating todo with title: %s", title)
				_, err := qn.service.CreateTodo(context.Background(), title, "")
				if err != nil {
					logger.Error("Failed to create todo: %v", err)
					qn.err = err
					return qn, tea.Quit
				}
				logger.Info("Quick note created: %s", title)
			}
			return qn, tea.Quit

		case tea.KeyEsc:
			logger.Debug("Escape key pressed")
			logger.Info("Quick note cancelled")
			return qn, tea.Quit

		case tea.KeyCtrlC:
			logger.Debug("Ctrl+C pressed")
			logger.Info("Quick note interrupted")
			return qn, tea.Quit
		}
	}

	qn.input, cmd = qn.input.Update(msg)
	return qn, cmd
}

func (qn *QuickNoteUI) View() string {
	if qn.err != nil {
		return "\n  Error: " + qn.err.Error() + "\n"
	}
	return "\n  " + qn.input.View() + "\n  (Enter to save, Esc to cancel)\n"
}
