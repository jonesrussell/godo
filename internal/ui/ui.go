package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/service"
)

// todoItem represents a todo in the list
type todoItem struct {
	todo model.Todo
}

// implement list.Item interface
func (i todoItem) Title() string       { return i.todo.Title }
func (i todoItem) Description() string { return i.todo.Description }
func (i todoItem) FilterValue() string { return i.todo.Title }

type TodoUI struct {
	todos    []model.Todo
	service  *service.TodoService
	cursor   int
	input    textinput.Model
	adding   bool
	err      error
}

func New(service *service.TodoService) *TodoUI {
	input := textinput.New()
	input.Placeholder = "Enter todo title..."
	input.Focus()

	return &TodoUI{
		service: service,
		input:   input,
		adding:  false,
	}
}

func (ui *TodoUI) Init() tea.Cmd {
	logger.Debug("Initializing TodoUI")
	return ui.loadTodos
}

type todosMsg struct {
	todos []model.Todo
	err   error
}

func (ui *TodoUI) loadTodos() tea.Msg {
	logger.Debug("Loading todos from service")
	todos, err := ui.service.ListTodos(context.Background())
	return todosMsg{todos: todos, err: err}
}

func (ui *TodoUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case todosMsg:
		if msg.err != nil {
			ui.err = msg.err
			return ui, nil
		}
		ui.todos = msg.todos
		if ui.cursor >= len(ui.todos) {
			ui.cursor = len(ui.todos) - 1
			if ui.cursor < 0 {
				ui.cursor = 0
			}
		}
		return ui, nil
	case ShowMsg:
		logger.SetUIActive(true)
		return ui, ui.loadTodos
	case tea.KeyMsg:
		if ui.adding {
			var cmd tea.Cmd
			ui.input, cmd = ui.input.Update(msg)
			
			if msg.String() == "enter" {
				title := ui.input.Value()
				if title != "" {
					_, err := ui.service.CreateTodo(context.Background(), title, "")
					if err != nil {
						ui.err = err
						ui.adding = false
						ui.input.Reset()
						return ui, nil
					}
				}
				ui.adding = false
				ui.input.Reset()
				return ui, ui.loadTodos
			}
			return ui, cmd
		}

		switch msg.String() {
		case "q":
			logger.SetUIActive(false)
			return ui, tea.Quit
		case "a":
			ui.adding = true
			ui.input.Focus()
			return ui, nil
		case "d":
			if len(ui.todos) > 0 && ui.cursor < len(ui.todos) {
				todoID := ui.todos[ui.cursor].ID
				if err := ui.service.DeleteTodo(context.Background(), todoID); err != nil {
					ui.err = err
					return ui, nil
				}
				return ui, ui.loadTodos
			}
		case " ":
			if len(ui.todos) > 0 && ui.cursor < len(ui.todos) {
				todoID := ui.todos[ui.cursor].ID
				if err := ui.service.ToggleTodoStatus(context.Background(), todoID); err != nil {
					ui.err = err
					return ui, nil
				}
				return ui, ui.loadTodos
			}
		case "up", "k":
			if ui.cursor > 0 {
				ui.cursor--
			}
		case "down", "j":
			if ui.cursor < len(ui.todos)-1 {
				ui.cursor++
			}
		}
	}
	
	return ui, nil
}

func (ui *TodoUI) View() string {
	var s strings.Builder

	if ui.adding {
		s.WriteString("\n  Add new todo:\n")
		s.WriteString(ui.input.View())
		s.WriteString("\n  (press enter to save)\n")
		return s.String()
	}

	// Show todos
	s.WriteString("\n  Todos:\n\n")
	
	if len(ui.todos) == 0 {
		s.WriteString("  No items\n")
	} else {
		for i, todo := range ui.todos {
			cursor := " "
			if ui.cursor == i {
				cursor = ">"
			}
			checkbox := "☐"
			if todo.Completed {
				checkbox = "☑"
			}
			s.WriteString(fmt.Sprintf("  %s %s %s\n", cursor, checkbox, todo.Title))
		}
	}

	s.WriteString("\n")
	
	// Help text
	if !ui.adding {
		s.WriteString("  a: add • d: delete • space: toggle • q: quit\n")
	}

	return s.String()
}

