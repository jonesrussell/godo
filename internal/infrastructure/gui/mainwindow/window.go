// Package mainwindow provides the main application window
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

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/shared/config"
)

// Window represents the main application window
type Window struct {
	app    fyne.App
	window fyne.Window
	store  storage.TaskStore
	log    logger.Logger
	tasks  []model.Task
	cfg    config.WindowConfig

	// UI components
	taskList    *widget.List
	addButton   *widget.Button
	refreshBtn  *widget.Button
	searchEntry *widget.Entry
	toolbar     *fyne.Container // Added for the new createMainLayout
}

// New creates a new main window
func New(app fyne.App, store storage.TaskStore, log logger.Logger, cfg config.WindowConfig) *Window {
	w := &Window{
		app:    app,
		store:  store,
		log:    log,
		cfg:    cfg,
		tasks:  make([]model.Task, 0),
		window: app.NewWindow("Godo - Task Manager"),
	}

	w.setupUI()
	w.loadTasks()
	return w
}

// Show displays the main window
func (w *Window) Show() {
	w.window.Show()
}

// Hide hides the main window
func (w *Window) Hide() {
	w.window.Hide()
}

// GetWindow returns the underlying Fyne window
func (w *Window) GetWindow() fyne.Window {
	return w.window
}

// SetContent sets the content of the main window
func (w *Window) SetContent(content fyne.CanvasObject) {
	w.window.SetContent(content)
}

// setupUI initializes the user interface
func (w *Window) setupUI() {
	w.createTaskList()
	w.createToolbar()
	w.createMainLayout()
}

// createTaskList creates the task list widget
func (w *Window) createTaskList() {
	w.taskList = widget.NewList(
		func() int { return len(w.tasks) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel("Task content"),
				layout.NewSpacer(),
				widget.NewButton("Edit", nil),
				widget.NewButton("Delete", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			box, ok := obj.(*fyne.Container)
			if !ok {
				w.log.Error("Failed to cast object to container")
				return
			}
			task := w.tasks[id]

			// Update check box
			check, ok := box.Objects[0].(*widget.Check)
			if !ok {
				w.log.Error("Failed to cast object to check")
				return
			}
			check.Checked = task.Done
			check.OnChanged = func(checked bool) {
				w.toggleTask(id, checked)
			}

			// Update label
			label, ok := box.Objects[1].(*widget.Label)
			if !ok {
				w.log.Error("Failed to cast object to label")
				return
			}
			label.SetText(task.Content)

			// Update edit button
			editBtn, ok := box.Objects[3].(*widget.Button)
			if !ok {
				w.log.Error("Failed to cast object to edit button")
				return
			}
			editBtn.OnTapped = func() {
				w.editTask(id)
			}

			// Update delete button
			deleteBtn, ok := box.Objects[4].(*widget.Button)
			if !ok {
				w.log.Error("Failed to cast object to delete button")
				return
			}
			deleteBtn.OnTapped = func() {
				w.deleteTask(id)
			}
		},
	)
}

// createToolbar creates the toolbar with buttons and search
func (w *Window) createToolbar() {
	// Create add button
	w.addButton = widget.NewButtonWithIcon("Add Task", theme.ContentAddIcon(), w.addTask)

	// Create refresh button
	w.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), w.loadTasks)

	// Create search entry
	w.searchEntry = widget.NewEntry()
	w.searchEntry.SetPlaceHolder("Search tasks...")
	w.searchEntry.OnChanged = w.filterTasks

	// Create toolbar
	w.toolbar = container.NewHBox(
		w.addButton,
		w.refreshBtn,
		layout.NewSpacer(),
		w.searchEntry,
	)
}

// createMainLayout creates the main window layout
func (w *Window) createMainLayout() {
	// Create main container
	content := container.NewBorder(
		w.toolbar,
		nil,
		nil,
		nil,
		w.taskList,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(600, 400))
}

// loadTasks loads all tasks from storage
func (w *Window) loadTasks() {
	ctx := context.Background()
	tasks, err := w.store.List(ctx)
	if err != nil {
		w.log.Error("Failed to load tasks", "error", err)
		return
	}

	w.tasks = tasks
	w.taskList.Refresh()
	w.log.Info("Tasks loaded", "count", len(tasks))
}

// addTask adds a new task
func (w *Window) addTask() {
	content := widget.NewEntry()
	content.SetPlaceHolder("Enter task content...")

	dialog := widget.NewForm(
		widget.NewFormItem("Task", content),
	)

	dialog.OnSubmit = func() {
		if content.Text == "" {
			return
		}

		task := model.Task{
			ID:        uuid.New().String(),
			Content:   content.Text,
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		ctx := context.Background()
		if err := w.store.Add(ctx, &task); err != nil {
			w.log.Error("Failed to add task", "error", err)
			return
		}

		w.loadTasks()
		dialog.Hide()
	}

	dialog.Show()
}

// toggleTask toggles the done status of a task
func (w *Window) toggleTask(id widget.ListItemID, done bool) {
	if id >= len(w.tasks) {
		return
	}

	task := w.tasks[id]
	task.Done = done
	task.UpdatedAt = time.Now()

	ctx := context.Background()
	if err := w.store.Update(ctx, &task); err != nil {
		w.log.Error("Failed to update task", "task_id", task.ID, "error", err)
		return
	}

	w.tasks[id] = task
	w.taskList.Refresh()
}

// editTask opens a dialog to edit a task
func (w *Window) editTask(id widget.ListItemID) {
	if id >= len(w.tasks) {
		return
	}

	task := w.tasks[id]
	content := widget.NewEntry()
	content.SetText(task.Content)

	dialog := widget.NewForm(
		widget.NewFormItem("Task", content),
	)

	dialog.OnSubmit = func() {
		if content.Text == "" {
			return
		}

		task.Content = content.Text
		task.UpdatedAt = time.Now()

		ctx := context.Background()
		if err := w.store.Update(ctx, &task); err != nil {
			w.log.Error("Failed to update task", "task_id", task.ID, "error", err)
			return
		}

		w.tasks[id] = task
		w.taskList.Refresh()
		dialog.Hide()
	}

	dialog.Show()
}

// deleteTask deletes a task
func (w *Window) deleteTask(id widget.ListItemID) {
	if id >= len(w.tasks) {
		return
	}

	task := w.tasks[id]
	ctx := context.Background()
	if err := w.store.Delete(ctx, task.ID); err != nil {
		w.log.Error("Failed to delete task", "task_id", task.ID, "error", err)
		return
	}

	w.loadTasks()
}

// filterTasks filters the task list based on search text
func (w *Window) filterTasks(searchText string) {
	// For now, just reload all tasks
	// In a real implementation, you might want to implement client-side filtering
	w.loadTasks()
}

// CenterOnScreen centers the main window on the screen
func (w *Window) CenterOnScreen() {
	w.window.CenterOnScreen()
}

// Refresh reloads and updates the task list
func (w *Window) Refresh() {
	w.loadTasks()
}

// Resize resizes the main window
func (w *Window) Resize(size fyne.Size) {
	w.window.Resize(size)
}
