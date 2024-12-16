package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	todoService *service.TodoService
	list        list.Model
	err         error
}

func New(todoService *service.TodoService) *TodoUI {
	logger.Debug("Creating new TodoUI instance")
	// Create a new list
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Todo List"
	l.SetShowHelp(true)

	return &TodoUI{
		todoService: todoService,
		list:        l,
	}
}

func (ui *TodoUI) Init() tea.Cmd {
	logger.Debug("Initializing TodoUI")
	// Initial command to load todos
	return ui.loadTodos
}

func (ui *TodoUI) loadTodos() tea.Msg {
	logger.Debug("Loading todos from service")
	todos, err := ui.todoService.ListTodos(context.Background())
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
		logger.Debug("Received key press: %s", msg.String())
		switch msg.String() {
		case "ctrl+c", "q":
			logger.Info("Quitting application")
			return ui, tea.Quit
		case "enter":
			if i, ok := ui.list.SelectedItem().(todoItem); ok {
				logger.Debug("Toggling todo status for ID: %d", i.todo.ID)
				// Toggle todo status
				if err := ui.todoService.ToggleTodoStatus(context.Background(), i.todo.ID); err != nil {
					logger.Error("Failed to toggle todo status: %v", err)
					ui.err = err
					return ui, nil
				}
				return ui, ui.loadTodos
			}
		}

	case errMsg:
		logger.Error("UI error: %v", msg.error)
		ui.err = msg
		return ui, nil

	case todosLoadedMsg:
		logger.Debug("Updating UI with loaded todos")
		ui.list.SetItems(msg.items)

	case ShowMsg:
		logger.Debug("Received show message")
		// Handle showing the UI
		return ui, ui.loadTodos
	}

	var cmd tea.Cmd
	ui.list, cmd = ui.list.Update(msg)
	return ui, cmd
}

func (ui *TodoUI) View() string {
	if ui.err != nil {
		logger.Error("Error in view: %v", ui.err)
		return fmt.Sprintf("Error: %v\n", ui.err)
	}
	return "\n" + ui.list.View()
}
