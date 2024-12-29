// Package mainwindow implements the main application window
package mainwindow

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the main window functionality
type Window struct {
	store  storage.TaskStore
	logger logger.Logger
	window fyne.Window
	config config.WindowConfig
}

// New creates a new main window
func New(app fyne.App, store storage.TaskStore, logger logger.Logger, cfg config.WindowConfig) *Window {
	window := app.NewWindow("Godo")
	w := &Window{
		store:  store,
		logger: logger,
		window: window,
		config: cfg,
	}

	// Set up window properties
	window.Resize(fyne.NewSize(float32(cfg.Width), float32(cfg.Height)))
	window.CenterOnScreen()
	window.SetIcon(theme.HomeIcon())

	// Create main content
	content := w.createContent()
	window.SetContent(content)

	// Set up window callbacks
	window.SetCloseIntercept(func() {
		w.Hide()
	})

	return w
}

// createContent creates the main window content
func (w *Window) createContent() fyne.CanvasObject {
	// Create task list
	taskList := widget.NewList(
		func() int { return 0 }, // TODO: Return actual task count
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel("Task content"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// TODO: Update item with actual task data
		},
	)

	// Create add task button
	addButton := widget.NewButtonWithIcon("Add Task", theme.ContentAddIcon(), func() {
		// TODO: Implement add task functionality
	})

	// Create toolbar
	toolbar := container.NewHBox(
		addButton,
	)

	// Create main layout
	return container.NewBorder(
		toolbar, nil, nil, nil, // top, bottom, left, right
		taskList, // center content
	)
}

// Show displays the window
func (w *Window) Show() {
	if !w.config.StartHidden {
		w.window.Show()
	}
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
