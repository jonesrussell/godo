// Package mainwindow implements the main application window
package mainwindow

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window represents the main application window
type Window struct {
	fyneWindow fyne.Window
	store      storage.Store
	logger     logger.Logger
	config     config.WindowConfig
	tasks      map[string]storage.Task
	list       *widget.List
}

// New creates a new main window instance
func New(app fyne.App, store storage.Store, logger logger.Logger, cfg config.WindowConfig) *Window {
	w := &Window{
		fyneWindow: app.NewWindow("Godo"),
		store:      store,
		logger:     logger,
		config:     cfg,
		tasks:      make(map[string]storage.Task),
	}

	w.setupUI()
	w.loadTasks()

	return w
}

// setupUI initializes the window UI components
func (w *Window) setupUI() {
	w.list = widget.NewList(
		func() int { return len(w.tasks) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			task := w.getTaskByIndex(id)
			if task == nil {
				return
			}

			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			check.Checked = task.Completed
			check.OnChanged = func(checked bool) {
				task.Completed = checked
				task.UpdatedAt = time.Now().Unix()
				w.updateTask(*task)
			}

			label.Text = task.Title
			label.Refresh()
		},
	)

	input := widget.NewEntry()
	input.SetPlaceHolder("Add a new task...")
	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}

		now := time.Now().Unix()
		task := storage.Task{
			ID:        fmt.Sprintf("%d", now), // Simple ID generation
			Title:     text,
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := w.store.Add(context.Background(), task); err != nil {
			w.logger.Error("Failed to add task", "error", err)
			return
		}

		w.tasks[task.ID] = task
		w.list.Refresh()
		input.SetText("")
	}

	content := container.NewBorder(input, nil, nil, nil, w.list)
	w.fyneWindow.SetContent(content)
	w.fyneWindow.Resize(fyne.NewSize(400, 300))
}

// loadTasks loads tasks from storage
func (w *Window) loadTasks() {
	tasks, err := w.store.List(context.Background())
	if err != nil {
		w.logger.Error("Failed to load tasks", "error", err)
		return
	}

	for _, task := range tasks {
		w.tasks[task.ID] = task
	}
	w.list.Refresh()
}

// updateTask updates a task in storage
func (w *Window) updateTask(task storage.Task) {
	if err := w.store.Update(context.Background(), task); err != nil {
		w.logger.Error("Failed to update task", "error", err)
		return
	}
	w.tasks[task.ID] = task
	w.list.Refresh()
}

// getTaskByIndex returns a task by its list index
func (w *Window) getTaskByIndex(index int) *storage.Task {
	i := 0
	for _, task := range w.tasks {
		if i == index {
			return &task
		}
		i++
	}
	return nil
}

// Show shows the window
func (w *Window) Show() {
	w.fyneWindow.Show()
}

// Hide hides the window
func (w *Window) Hide() {
	w.fyneWindow.Hide()
}

// Close closes the window
func (w *Window) Close() {
	w.fyneWindow.Close()
}

// RequestFocus requests focus for the window
func (w *Window) RequestFocus() {
	w.fyneWindow.RequestFocus()
}

// CenterOnScreen centers the window on screen
func (w *Window) CenterOnScreen() {
	w.fyneWindow.CenterOnScreen()
}

// Resize resizes the window
func (w *Window) Resize(size fyne.Size) {
	w.fyneWindow.Resize(size)
}

// SetOnClosed sets the window close callback
func (w *Window) SetOnClosed(callback func()) {
	w.fyneWindow.SetOnClosed(callback)
}

// SetCloseIntercept sets the window close intercept callback
func (w *Window) SetCloseIntercept(callback func()) {
	w.fyneWindow.SetCloseIntercept(callback)
}

// Canvas returns the window canvas
func (w *Window) Canvas() fyne.Canvas {
	return w.fyneWindow.Canvas()
}
