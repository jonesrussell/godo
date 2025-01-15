// Package quicknote implements the quick note window functionality
package quicknote

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the quick note window
type Window struct {
	store      storage.TaskStore
	logger     logger.Logger
	window     fyne.Window
	app        fyne.App
	config     config.WindowConfig
	input      *Entry
	saveBtn    *widget.Button
	mainWindow mainwindow.Interface
}

// New creates a new quick note window
func New(app fyne.App, store storage.TaskStore, logger logger.Logger, config config.WindowConfig, mainWindow mainwindow.Interface) *Window {
	w := &Window{
		store:      store,
		logger:     logger,
		app:        app,
		config:     config,
		window:     app.NewWindow("Quick Note"),
		mainWindow: mainWindow,
	}

	w.setupUI()
	return w
}

// setupUI initializes the window's UI components
func (w *Window) setupUI() {
	w.input = NewEntry()
	w.input.SetPlaceHolder("Enter your quick note...")
	w.input.SetOnCtrlEnter(w.saveNote)

	w.saveBtn = widget.NewButton("Save", w.saveNote)

	content := container.NewBorder(
		nil,       // top
		w.saveBtn, // bottom
		nil,       // left
		nil,       // right
		w.input,   // center
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(float32(w.config.Width), float32(w.config.Height)))
	w.window.CenterOnScreen()
}

// saveNote saves the current note
func (w *Window) saveNote() {
	if w.input.Text != "" {
		now := time.Now()
		task := storage.Task{
			ID:        uuid.New().String(),
			Content:   w.input.Text,
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := w.store.Add(context.Background(), task); err != nil {
			w.logger.Error("Failed to add quick note", "error", err)
			return
		}
		w.input.SetText("")
		w.Hide()

		// Refresh the main window
		if w.mainWindow != nil {
			w.mainWindow.Refresh()
		}
	}
}

// Show displays the quick note window
func (w *Window) Show() {
	w.input.SetText("")
	w.window.Show()
	w.window.Canvas().Focus(w.input)
}

// Hide hides the quick note window
func (w *Window) Hide() {
	w.window.Hide()
}
