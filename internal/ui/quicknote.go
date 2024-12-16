package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/service"
)

type QuickNoteUI struct {
	input   textinput.Model
	service service.TodoServicer
	err     error
}

func NewQuickNote(service service.TodoServicer) *QuickNoteUI {
	input := textinput.New()
	input.Placeholder = "Type your todo and press Enter..."
	input.Focus()
	input.Width = 50

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
			if text := qn.input.Value(); text != "" {
				_, err := qn.service.CreateTodo(context.Background(), text, "")
				if err != nil {
					qn.err = err
					return qn, tea.Quit
				}
				return qn, tea.Quit
			}
			return qn, tea.Quit

		case tea.KeyEsc:
			return qn, tea.Quit

		case tea.KeyCtrlC:
			return qn, tea.Quit
		}
	}

	qn.input, cmd = qn.input.Update(msg)
	return qn, cmd
}

func (qn *QuickNoteUI) View() string {
	if qn.err != nil {
		return fmt.Sprintf("\n  Error: %v\n\n", qn.err)
	}
	return "\n  " + qn.input.View() + "\n\n"
}
