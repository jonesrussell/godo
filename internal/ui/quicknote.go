package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
)

type QuickNote struct {
	textInput textinput.Model
	service   *service.TodoService
	err       error
}

func NewQuickNote(service *service.TodoService) *QuickNote {
	ti := textinput.New()
	ti.Placeholder = "Enter quick note..."
	ti.Focus()
	ti.Width = 40

	return &QuickNote{
		textInput: ti,
		service:   service,
	}
}

func (qn *QuickNote) Init() tea.Cmd {
	logger.Debug("Initializing QuickNote UI")
	return textinput.Blink
}

func (qn *QuickNote) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if text := qn.textInput.Value(); text != "" {
				todo, err := qn.service.CreateTodo(context.Background(), "quick", text)
				if err != nil {
					logger.Error("Failed to create todo: %v", err)
					qn.err = err
					return qn, tea.Quit
				}
				logger.Debug("Created quick note: %s (ID: %d)", text, todo.ID)
				return qn, tea.Quit
			}
			return qn, nil
		case tea.KeyEsc:
			logger.Debug("Quick note cancelled")
			return qn, tea.Quit
		}
	}

	qn.textInput, cmd = qn.textInput.Update(msg)
	return qn, cmd
}

func (qn *QuickNote) View() string {
	if qn.err != nil {
		return "Error: " + qn.err.Error() + "\nPress any key to exit."
	}

	return "\n  🗒️  Quick Note\n\n" +
		"  " + qn.textInput.View() + "\n\n" +
		"  (Enter to save • Esc to cancel)\n"
}
