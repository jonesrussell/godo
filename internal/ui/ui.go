package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	// Initial command to load todos
	return ui.loadTodos
}

func (ui *TodoUI) loadTodos() tea.Msg {
	todos, err := ui.todoService.ListTodos(context.Background())
	if err != nil {
		return errMsg{err}
	}

	items := make([]list.Item, len(todos))
	for i, todo := range todos {
		items[i] = todoItem{todo: todo}
	}
	return todosLoadedMsg{items}
}

// Custom messages
type errMsg struct{ error }
type todosLoadedMsg struct{ items []list.Item }

func (ui *TodoUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return ui, tea.Quit
		case "enter":
			if i, ok := ui.list.SelectedItem().(todoItem); ok {
				// Toggle todo status
				if err := ui.todoService.ToggleTodoStatus(context.Background(), i.todo.ID); err != nil {
					ui.err = err
					return ui, nil
				}
				return ui, ui.loadTodos
			}
		}

	case errMsg:
		ui.err = msg
		return ui, nil

	case todosLoadedMsg:
		ui.list.SetItems(msg.items)
	}

	var cmd tea.Cmd
	ui.list, cmd = ui.list.Update(msg)
	return ui, cmd
}

func (ui *TodoUI) View() string {
	if ui.err != nil {
		return fmt.Sprintf("Error: %v\n", ui.err)
	}
	return "\n" + ui.list.View()
}
