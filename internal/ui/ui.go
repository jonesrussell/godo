package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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
	// Initial command to load todos
	return ui.loadTodos
}

func (ui *TodoUI) loadTodos() tea.Msg {
	logger.Debug("Loading todos from service")
	todos, err := ui.service.ListTodos(context.Background())
	if err != nil {
		logger.Error("Failed to load todos: %v", err)
		return errMsg{err}
	}

	items := make([]list.Item, len(todos))
	for i, todo := range todos {
		items[i] = todoItem{todo: todo}
	}
	logger.Debug("Loaded %d todos", len(todos))
	return todosLoadedMsg{items}
}

// Custom messages
type errMsg struct{ error }
type todosLoadedMsg struct{ items []list.Item }

func (ui *TodoUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			ui.adding = true
			ui.input.Focus()
			return ui, nil
		case "enter":
			if ui.adding {
				title := ui.input.Value()
				if err := ui.service.CreateTodo(context.Background(), title, ""); err != nil {
					ui.err = err
				}
				ui.adding = false
				ui.input.Reset()
				return ui, ui.loadTodos
			}
		}
	}
	
	if ui.adding {
		var cmd tea.Cmd
		ui.input, cmd = ui.input.Update(msg)
		return ui, cmd
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

	if ui.err != nil {
		logger.Error("Error in view: %v", ui.err)
		return fmt.Sprintf("Error: %v\n", ui.err)
	}
	return "\n" + ui.list.View()
}
