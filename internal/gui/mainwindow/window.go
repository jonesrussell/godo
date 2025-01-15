// Package mainwindow implements the main application window
package mainwindow

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the main window functionality
type Window struct {
	store  storage.TaskStore
	logger logger.Logger
	window fyne.Window
	app    fyne.App
	config config.WindowConfig
	tasks  []storage.Task
}

// New creates a new main window
func New(app fyne.App, store storage.TaskStore, logger logger.Logger, config config.WindowConfig) *Window {
	w := &Window{
		store:  store,
		logger: logger,
		app:    app,
		config: config,
		window: app.NewWindow("Godo"),
	}

	w.setupUI()
	return w
}

// setupUI initializes the window's UI components
func (w *Window) setupUI() {
	// Create task list
	var err error
	w.tasks, err = w.store.List(context.Background())
	if err != nil {
		w.logger.Error("Failed to load tasks", "error", err)
	}

	taskList := widget.NewList(
		func() int { return len(w.tasks) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			check.Checked = w.tasks[id].Done
			check.OnChanged = func(done bool) {
				w.tasks[id].Done = done
				w.tasks[id].UpdatedAt = time.Now()
				if err := w.store.Update(context.Background(), w.tasks[id]); err != nil {
					w.logger.Error("Failed to update task", "error", err)
				}
			}
			label.SetText(w.tasks[id].Content)
		},
	)

	// Create add task input
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter a new task...")

	addButton := widget.NewButton("Add", func() {
		if input.Text != "" {
			now := time.Now()
			task := storage.Task{
				ID:        uuid.New().String(),
				Content:   input.Text,
				Done:      false,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := w.store.Add(context.Background(), task); err != nil {
				w.logger.Error("Failed to add task", "error", err)
				return
			}
			w.tasks = append(w.tasks, task)
			taskList.Refresh()
			input.SetText("")
		}
	})

	// Layout
	content := container.NewBorder(
		container.NewBorder(nil, nil, nil, addButton, input), // top
		nil, // bottom
		nil, // left
		nil, // right
		taskList,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(float32(w.config.Width), float32(w.config.Height)))
	w.window.CenterOnScreen()
}

// Show displays the window
func (w *Window) Show() {
	w.window.Show()
}

// Hide hides the window
func (w *Window) Hide() {
	w.window.Hide()
}

// SetContent sets the window's content
func (w *Window) SetContent(content fyne.CanvasObject) {
	w.window.SetContent(content)
}

// Resize changes the window's size
func (w *Window) Resize(size fyne.Size) {
	w.window.Resize(size)
}

// CenterOnScreen centers the window on the screen
func (w *Window) CenterOnScreen() {
	w.window.CenterOnScreen()
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.window
}

// Refresh reloads and updates the task list
func (w *Window) Refresh() {
	tasks, err := w.store.List(context.Background())
	if err != nil {
		w.logger.Error("Failed to reload tasks", "error", err)
		return
	}

	w.tasks = tasks
	if content, ok := w.window.Content().(*fyne.Container); ok {
		if list, ok := content.Objects[0].(*widget.List); ok {
			list.Refresh()
		}
	}
}
