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
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// QuickNote implements the quick note window functionality
type QuickNote struct {
	store  storage.TaskStore
	logger logger.Logger
	window fyne.Window
	app    fyne.App
	config config.WindowConfig
	input  *widget.Entry
}

// NewQuickNote creates a new quick note window
func NewQuickNote(app fyne.App, store storage.TaskStore, logger logger.Logger, config config.WindowConfig) *QuickNote {
	w := &QuickNote{
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
func (w *QuickNote) setupUI() {
	w.input = widget.NewMultiLineEntry()
	w.input.SetPlaceHolder("Enter your quick note...")

	saveButton := widget.NewButton("Save", func() {
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
				w.logger.Error("Failed to save quick note", "error", err)
				return
			}
			w.input.SetText("")
			w.Hide()
		}
	})

	content := container.NewBorder(
		nil,        // top
		saveButton, // bottom
		nil,        // left
		nil,        // right
		w.input,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(float32(w.config.Width), float32(w.config.Height)))
	w.window.CenterOnScreen()
}

// Show displays the window and focuses the input
func (w *QuickNote) Show() {
	w.input.SetText("")
	w.window.Show()
	w.window.Canvas().Focus(w.input)
}

// Hide hides the window
func (w *QuickNote) Hide() {
	w.window.Hide()
}

// SetContent sets the window's content
func (w *QuickNote) SetContent(content fyne.CanvasObject) {
	w.window.SetContent(content)
}

// Resize changes the window's size
func (w *QuickNote) Resize(size fyne.Size) {
	w.window.Resize(size)
}

// CenterOnScreen centers the window on the screen
func (w *QuickNote) CenterOnScreen() {
	w.window.CenterOnScreen()
}

// GetWindow returns the underlying fyne.Window
func (w *QuickNote) GetWindow() fyne.Window {
	return w.window
}
