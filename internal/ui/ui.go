package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
)

type TodoUI struct {
	todoService service.TodoService
	window      fyne.Window
}

func NewTodoUI(todoService service.TodoService, window fyne.Window) *TodoUI {
	return &TodoUI{
		todoService: todoService,
		window:      window,
	}
}

func (ui *TodoUI) Run() error {
	var taskList *widget.List

	// Create main layout
	taskList = widget.NewList(
		// Length function returns the total number of items
		func() int {
			todos, err := ui.todoService.ListTodos(context.Background())
			if err != nil {
				logger.Error("Failed to get todos", "error", err)
				return 0
			}
			return len(todos)
		},
		// CreateItem function returns a new template item
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
				widget.NewButton("Delete", nil),
			)
		},
		// UpdateItem function updates an item with real data
		func(id widget.ListItemID, item fyne.CanvasObject) {
			todos, err := ui.todoService.ListTodos(context.Background())
			if err != nil {
				logger.Error("Failed to get todos", "error", err)
				return
			}
			if id >= len(todos) {
				return
			}

			todo := todos[id]
			box := item.(*fyne.Container)

			// Update checkbox
			check := box.Objects[0].(*widget.Check)
			check.Checked = todo.Completed
			check.OnChanged = func(checked bool) {
				if err := ui.todoService.ToggleTodoStatus(context.Background(), todo.ID); err != nil {
					logger.Error("Failed to toggle todo status",
						"id", todo.ID,
						"error", err)
				}
				taskList.Refresh()
			}

			// Update label
			label := box.Objects[1].(*widget.Label)
			label.SetText(todo.Title)

			// Update delete button
			deleteBtn := box.Objects[2].(*widget.Button)
			deleteBtn.OnTapped = func() {
				if err := ui.todoService.DeleteTodo(context.Background(), todo.ID); err != nil {
					logger.Error("Failed to delete todo",
						"id", todo.ID,
						"error", err)
				}
				taskList.Refresh()
			}
		},
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("Add new task...")
	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}

		_, err := ui.todoService.CreateTodo(context.Background(), text, "")
		if err != nil {
			logger.Error("Failed to create todo",
				"title", text,
				"error", err)
			return
		}

		input.SetText("")
		taskList.Refresh()
	}

	content := container.NewBorder(input, nil, nil, nil, taskList)
	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(300, 500))
	ui.window.ShowAndRun()

	return nil
}
