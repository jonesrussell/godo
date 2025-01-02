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
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Window implements the QuickNoteManager interface
type Window struct {
	store   storage.TaskStore
	logger  logger.Logger
	window  fyne.Window
	app     fyne.App
	config  config.WindowConfig
	input   *widget.Entry
	saveBtn *widget.Button
}

// Ensure Window implements QuickNoteManager
var _ gui.QuickNoteManager = (*Window)(nil)

// New creates a new quick note window
func New(app fyne.App, store storage.TaskStore, logger logger.Logger, config config.WindowConfig) *Window {
	w := &Window{
		store:  store,
		logger: logger,
		app:    app,
		config: config,
		window: app.NewWindow("Quick Note"),
	}

	w.setupUI()
	return w
}

// setupUI initializes the window's UI components
func (w *Window) setupUI() {
	w.input = widget.NewMultiLineEntry()
	w.input.SetPlaceHolder("Enter your quick note...")

	w.saveBtn = widget.NewButton("Save", func() {
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
		}
	})

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

// Show displays the window
func (w *Window) Show() {
	w.window.Show()
	if w.input != nil {
		w.window.Canvas().Focus(w.input)
	}
}

// Hide hides the window
func (w *Window) Hide() {
	w.window.Hide()
}

// CenterOnScreen centers the window on the screen
func (w *Window) CenterOnScreen() {
	w.window.CenterOnScreen()
}
