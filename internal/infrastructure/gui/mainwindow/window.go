// Package mainwindow provides the main application window
package mainwindow

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/gui/windowmanager"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
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
	toolbar     *fyne.Container
	statusBar   *widget.Label

	// Focus management
	focusManager *windowmanager.FocusManager
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

	w.focusManager = windowmanager.NewFocusManager(w.window, log)
	w.setupUI()
	w.loadTasks()
	return w
}

// Show displays the main window
func (w *Window) Show() {
	w.window.Show()
	w.focusManager.RequestFocus()
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
	w.createStatusBar()
	w.createMainLayout()
	w.buildFocusQueue()
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
			w.updateTaskListItem(id, obj)
		},
	)

	// Add task list to focus queue
	w.focusManager.AddToFocusQueue(w.taskList)
}

// updateTaskListItem updates a task list item
func (w *Window) updateTaskListItem(id widget.ListItemID, obj fyne.CanvasObject) {
	box, ok := obj.(*fyne.Container)
	if !ok {
		w.log.Error("Failed to cast object to container")
		return
	}

	if id >= len(w.tasks) {
		return
	}
	task := w.tasks[id]

	// Update check box
	if check, ok := box.Objects[0].(*widget.Check); ok {
		check.Checked = task.Done
		check.OnChanged = func(checked bool) {
			w.toggleTask(id, checked)
		}
		w.focusManager.AddToFocusQueue(check)
	}

	// Update label
	if label, ok := box.Objects[1].(*widget.Label); ok {
		label.SetText(task.Content)
	}

	// Update edit button
	if editBtn, ok := box.Objects[3].(*widget.Button); ok {
		editBtn.OnTapped = func() {
			w.editTask(id)
		}
		w.focusManager.AddToFocusQueue(editBtn)
	}

	// Update delete button
	if deleteBtn, ok := box.Objects[4].(*widget.Button); ok {
		deleteBtn.OnTapped = func() {
			w.deleteTask(id)
		}
		w.focusManager.AddToFocusQueue(deleteBtn)
	}
}

// createToolbar creates the toolbar with buttons and search
func (w *Window) createToolbar() {
	// Create add button
	w.addButton = widget.NewButtonWithIcon("Add Task", theme.ContentAddIcon(), w.addTask)
	w.focusManager.AddToFocusQueue(w.addButton)

	// Create refresh button
	w.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), w.loadTasks)
	w.focusManager.AddToFocusQueue(w.refreshBtn)

	// Create search entry
	w.searchEntry = widget.NewEntry()
	w.searchEntry.SetPlaceHolder("Search tasks...")
	w.searchEntry.OnChanged = w.filterTasks
	w.focusManager.AddToFocusQueue(w.searchEntry)

	// Create toolbar
	w.toolbar = container.NewHBox(
		w.addButton,
		w.refreshBtn,
		layout.NewSpacer(),
		w.searchEntry,
	)
}

// createStatusBar creates the status bar
func (w *Window) createStatusBar() {
	w.statusBar = widget.NewLabel("")
	w.statusBar.Hide()
}

// createMainLayout creates the main window layout
func (w *Window) createMainLayout() {
	// Create main container
	content := container.NewBorder(
		w.toolbar,
		w.statusBar,
		nil,
		nil,
		w.taskList,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(600, 400))
}

// buildFocusQueue builds the focus queue for keyboard navigation
func (w *Window) buildFocusQueue() {
	w.focusManager.BuildFocusQueue()
}

// loadTasks loads all tasks from storage
func (w *Window) loadTasks() {
	ctx := context.Background()
	tasks, err := w.store.List(ctx)
	if err != nil {
		w.log.Error("Failed to load tasks", "error", err)
		w.showStatus("Failed to load tasks", true)
		return
	}

	w.tasks = tasks
	w.taskList.Refresh()
	w.showStatus(fmt.Sprintf("Loaded %d tasks", len(tasks)), false)
	w.log.Info("Tasks loaded", "count", len(tasks))
}

// addTask adds a new task with enhanced dialog
func (w *Window) addTask() {
	content := widget.NewEntry()
	content.SetPlaceHolder("Enter task content...")
	w.focusManager.AddToFocusQueue(content)

	form := dialog.NewForm(
		"Add Task",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Task", content),
		},
		func(confirm bool) {
			if !confirm || content.Text == "" {
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
				w.showStatus("Failed to add task", true)
				return
			}

			w.loadTasks()
			w.showStatus("Task added successfully", false)
		},
		w.window,
	)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
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
		w.showStatus("Failed to update task", true)
		return
	}

	w.tasks[id] = task
	w.taskList.Refresh()
	w.showStatus("Task updated", false)
}

// editTask opens a dialog to edit a task
func (w *Window) editTask(id widget.ListItemID) {
	if id >= len(w.tasks) {
		return
	}

	task := w.tasks[id]
	content := widget.NewEntry()
	content.SetText(task.Content)
	w.focusManager.AddToFocusQueue(content)

	form := dialog.NewForm(
		"Edit Task",
		"Save",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Task", content),
		},
		func(confirm bool) {
			if !confirm || content.Text == "" {
				return
			}

			task.Content = content.Text
			task.UpdatedAt = time.Now()

			ctx := context.Background()
			if err := w.store.Update(ctx, &task); err != nil {
				w.log.Error("Failed to update task", "task_id", task.ID, "error", err)
				w.showStatus("Failed to update task", true)
				return
			}

			w.tasks[id] = task
			w.taskList.Refresh()
			w.showStatus("Task updated", false)
		},
		w.window,
	)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
}

// deleteTask deletes a task with confirmation
func (w *Window) deleteTask(id widget.ListItemID) {
	if id >= len(w.tasks) {
		return
	}

	task := w.tasks[id]

	confirm := dialog.NewConfirm(
		"Delete Task",
		fmt.Sprintf("Are you sure you want to delete '%s'?", task.Content),
		func(confirm bool) {
			if !confirm {
				return
			}

			ctx := context.Background()
			if err := w.store.Delete(ctx, task.ID); err != nil {
				w.log.Error("Failed to delete task", "task_id", task.ID, "error", err)
				w.showStatus("Failed to delete task", true)
				return
			}

			w.loadTasks()
			w.showStatus("Task deleted", false)
		},
		w.window,
	)

	confirm.Show()
}

// filterTasks filters the task list based on search text
func (w *Window) filterTasks(searchText string) {
	// For now, just reload all tasks
	// In a real implementation, you might want to implement client-side filtering
	w.loadTasks()
}

// showStatus shows a status message
func (w *Window) showStatus(message string, isError bool) {
	fyne.Do(func() {
		w.statusBar.SetText(message)
		w.statusBar.Show()

		if isError {
			w.statusBar.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			w.statusBar.TextStyle = fyne.TextStyle{}
		}
	})

	// Auto-hide after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		fyne.Do(func() {
			w.statusBar.Hide()
		})
	}()
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
