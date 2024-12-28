package mainwindow

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	store    storage.Store
	logger   logger.Logger
	win      fyne.Window
	tasks    []storage.Task
	taskList *widget.List
}

// New creates a new main window
func New(store storage.Store, logger logger.Logger) *Window {
	return &Window{
		store:  store,
		logger: logger,
	}
}

// Setup initializes the window
func (w *Window) Setup() error {
	fmt.Println("Starting main window setup")
	w.logger.Debug("Starting main window setup")

	// Basic implementation to use the win field
	fmt.Println("Creating new window")
	w.logger.Debug("Creating new window")
	w.win = fyne.CurrentApp().NewWindow("Godo")
	if w.win == nil {
		fmt.Println("ERROR: Failed to create window - NewWindow returned nil")
		w.logger.Error("Failed to create window - NewWindow returned nil")
		return fmt.Errorf("failed to create window")
	}
	fmt.Println("Window created successfully")
	w.logger.Debug("Window created successfully")

	fmt.Println("Setting window properties")
	w.win.SetMaster() // Make this window the main window
	w.win.SetCloseIntercept(func() {
		fmt.Println("Window close requested")
		w.win.Hide()
	})

	fmt.Println("Resizing window")
	w.logger.Debug("Resizing window", "width", 800, "height", 600)
	w.win.Resize(fyne.NewSize(800, 600))
	w.win.CenterOnScreen()

	fmt.Println("Creating main content")
	w.logger.Debug("Creating main content")

	// Load initial tasks
	tasks, err := w.store.List()
	if err != nil {
		w.logger.Error("Failed to load tasks", "error", err)
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	w.tasks = tasks

	// Create task list
	w.taskList = widget.NewList(
		func() int {
			return len(w.tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
				widget.NewButton("Delete", nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			task := w.tasks[id]
			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)
			deleteBtn := box.Objects[2].(*widget.Button)

			// Update check state
			check.Checked = task.Completed
			check.OnChanged = func(checked bool) {
				task.Completed = checked
				if err := w.store.Update(task); err != nil {
					w.logger.Error("Failed to update task", "error", err)
					// Revert UI state on error
					check.Checked = !checked
					check.Refresh()
					return
				}
				w.tasks[id] = task
			}

			// Update label
			label.SetText(task.Title)

			// Setup delete button
			deleteBtn.OnTapped = func() {
				if err := w.store.Delete(task.ID); err != nil {
					w.logger.Error("Failed to delete task", "error", err)
					return
				}
				// Remove from local list and refresh
				w.tasks = append(w.tasks[:id], w.tasks[id+1:]...)
				w.taskList.Refresh()
			}
		},
	)

	// Create refresh button
	refreshBtn := widget.NewButton("Refresh", func() {
		tasks, err := w.store.List()
		if err != nil {
			w.logger.Error("Failed to refresh tasks", "error", err)
			return
		}
		w.tasks = tasks
		w.taskList.Refresh()
	})

	content := container.NewBorder(
		container.NewHBox(
			widget.NewLabel("Tasks"),
			refreshBtn,
		),
		nil, nil, nil,
		w.taskList,
	)
	w.win.SetContent(content)

	// Don't show the window on startup
	w.win.Hide()

	fmt.Println("Window setup complete")
	w.logger.Debug("Window setup complete")

	return nil
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.win
}

// Show displays the window
func (w *Window) Show() {
	if w.win != nil {
		w.win.Show()
	} else {
		w.logger.Error("Cannot show window: window not initialized")
	}
}

// Hide hides the window
func (w *Window) Hide() {
	if w.win != nil {
		w.win.Hide()
	} else {
		w.logger.Error("Cannot hide window: window not initialized")
	}
}
