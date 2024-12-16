package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
)

type QuickNoteUI struct {
	input   textinput.Model
	service *service.TodoService
	err     error
}

func NewQuickNote(service *service.TodoService) *QuickNoteUI {
	input := textinput.New()
	input.Placeholder = "Type your note and press Enter..."
	input.Focus()
	input.Width = 50 // Adjust as needed

	return &QuickNoteUI{
		input:   input,
		service: service,
	}
}

func (qn *QuickNoteUI) Init() tea.Cmd {
	return textinput.Blink
}

func (qn *QuickNoteUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if title := qn.input.Value(); title != "" {
				_, err := qn.service.CreateTodo(context.Background(), title, "")
				if err != nil {
					logger.Error("Failed to create todo", err)
					qn.err = err
					return qn, tea.Quit
				}
				logger.Info("Quick note created: " + title)
			}
			return qn, tea.Quit

		case tea.KeyEsc:
			logger.Info("Quick note cancelled")
			return qn, tea.Quit

		case tea.KeyCtrlC:
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
