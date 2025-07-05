// Package mainwindow implements the main application window
package mainwindow

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// TaskManager handles task-related operations
type TaskManager struct {
	store  storage.TaskStore
	logger logger.Logger
	tasks  []storage.Task
}

// Window implements the main window functionality
type Window struct {
	*TaskManager
	window fyne.Window
	app    fyne.App
	config config.WindowConfig
	list   *widget.List
}

// New creates a new main window
func New(app fyne.App, store storage.TaskStore, logger logger.Logger, config config.WindowConfig) *Window {
	w := &Window{
		TaskManager: &TaskManager{
			store:  store,
			logger: logger,
			tasks:  make([]storage.Task, 0),
		},
		app:    app,
		config: config,
		window: app.NewWindow("Godo"),
	}

	if err := w.loadTasks(context.Background()); err != nil {
		w.logger.Error("Failed to load tasks during initialization", "error", err)
	}
	w.setupUI()
	return w
}

// loadTasks loads tasks from storage
func (tm *TaskManager) loadTasks(ctx context.Context) error {
	tasks, err := tm.store.List(ctx)
	if err != nil {
		tm.logger.Error("Failed to load tasks", "error", err)
		return err
	}
	tm.tasks = tasks
	return nil
}

// addTask adds a new task
func (tm *TaskManager) addTask(ctx context.Context, content string) error {
	if content == "" {
		return nil
	}

	task := storage.Task{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tm.store.Add(ctx, task); err != nil {
		tm.logger.Error("Failed to add task", "error", err)
		return err
	}

	tm.tasks = append(tm.tasks, task)
	return nil
}

// updateTask updates a task's status
func (tm *TaskManager) updateTask(ctx context.Context, id int, done bool) error {
	if id < 0 || id >= len(tm.tasks) {
		return nil
	}

	tm.tasks[id].Done = done
	tm.tasks[id].UpdatedAt = time.Now()

	if err := tm.store.Update(ctx, tm.tasks[id]); err != nil {
		tm.logger.Error("Failed to update task", "error", err)
		return err
	}
	return nil
}

// setupUI initializes the window's UI components
func (w *Window) setupUI() {
	w.list = w.createTaskList()
	input := w.createInput()
	addButton := w.createAddButton(input)
	content := w.createLayout(input, addButton)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(float32(w.config.Width), float32(w.config.Height)))
	w.window.CenterOnScreen()

	// Set window title and icon
	w.window.SetTitle("Godo - Task Manager")
	icon := theme.ListIcon()
	w.window.SetIcon(icon)

	// Focus input by default for better UX
	w.window.Canvas().Focus(input)
}

func (w *Window) createTaskList() *widget.List {
	return widget.NewList(
		func() int { return len(w.tasks) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
				layout.NewSpacer(),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			box, ok := item.(*fyne.Container)
			if !ok {
				w.logger.Error("Failed to cast item to container")
				return
			}
			check, ok := box.Objects[0].(*widget.Check)
			if !ok {
				w.logger.Error("Failed to cast object to check")
				return
			}
			label, ok := box.Objects[1].(*widget.Label)
			if !ok {
				w.logger.Error("Failed to cast object to label")
				return
			}

			check.Checked = w.tasks[id].Done
			check.OnChanged = func(done bool) {
				if err := w.updateTask(context.Background(), id, done); err != nil {
					w.logger.Error("Failed to update task", "error", err)
				}
			}
			label.SetText(w.tasks[id].Content)
			// Add visual feedback for completed tasks
			if w.tasks[id].Done {
				label.TextStyle = fyne.TextStyle{Monospace: true}
			} else {
				label.TextStyle = fyne.TextStyle{}
			}
		},
	)
}

func (w *Window) createInput() *widget.Entry {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter a new task...")
	// Add keyboard shortcut for quick task entry
	input.OnSubmitted = func(text string) {
		if err := w.addTask(context.Background(), text); err != nil {
			w.logger.Error("Failed to add task", "error", err)
		} else {
			input.SetText("")
			w.list.Refresh()
		}
	}
	return input
}

func (w *Window) createAddButton(input *widget.Entry) *widget.Button {
	return widget.NewButton("Add", func() {
		if err := w.addTask(context.Background(), input.Text); err != nil {
			w.logger.Error("Failed to add task", "error", err)
		} else {
			input.SetText("")
			w.list.Refresh()
		}
	})
}

func (w *Window) createLayout(input *widget.Entry, addButton *widget.Button) fyne.CanvasObject {
	// Create a more structured layout with padding
	inputContainer := container.NewBorder(nil, nil, nil, addButton, input)
	content := container.NewBorder(
		container.NewPadded(inputContainer), // Add padding around input
		nil,
		nil,
		nil,
		container.NewPadded(w.list), // Add padding around list
	)
	return content
}

// Window interface methods
func (w *Window) Show()                                { w.window.Show() }
func (w *Window) Hide()                                { w.window.Hide() }
func (w *Window) SetContent(content fyne.CanvasObject) { w.window.SetContent(content) }
func (w *Window) Resize(size fyne.Size)                { w.window.Resize(size) }
func (w *Window) CenterOnScreen()                      { w.window.CenterOnScreen() }
func (w *Window) GetWindow() fyne.Window               { return w.window }

// Refresh reloads and updates the task list
func (w *Window) Refresh() {
	if err := w.loadTasks(context.Background()); err != nil {
		w.logger.Error("Failed to load tasks during refresh", "error", err)
		return
	}
	if w.list != nil {
		w.list.Refresh()
	}
}
