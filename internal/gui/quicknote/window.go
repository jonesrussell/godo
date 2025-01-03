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
	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// Window represents the quick note window
type Window struct {
	fyneWindow fyne.Window
	store      types.Store
	logger     logger.Logger
	config     config.WindowConfig
	input      *widget.Entry
	saveBtn    *widget.Button
}

// New creates a new quick note window instance
func New(app fyne.App, store types.Store, logger logger.Logger, cfg config.WindowConfig) *Window {
	w := &Window{
		fyneWindow: app.NewWindow("Quick Note"),
		store:      store,
		logger:     logger,
		config:     cfg,
	}

	w.setupUI()
	return w
}

// setupUI initializes the window UI components
func (w *Window) setupUI() {
	w.input = widget.NewMultiLineEntry()
	w.input.SetPlaceHolder("Enter your quick note...")

	w.saveBtn = widget.NewButton("Save", func() {
		w.saveNote()
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		w.input.SetText("")
		w.Hide()
	})

	buttons := container.NewHBox(w.saveBtn, cancelBtn)
	content := container.NewBorder(nil, buttons, nil, nil, w.input)

	w.fyneWindow.SetContent(content)
	w.fyneWindow.Resize(fyne.NewSize(300, 200))
	w.fyneWindow.SetFixedSize(true)
}

// saveNote saves the current note text
func (w *Window) saveNote() {
	text := w.input.Text
	if text == "" {
		return
	}

	note := &note.Note{
		ID:        uuid.New().String(),
		Content:   text,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := w.store.Add(context.Background(), note); err != nil {
		w.logger.Error("Failed to save note", "error", err)
		return
	}

	w.input.SetText("")
	w.Hide()
}

// Show shows the window
func (w *Window) Show() {
	w.fyneWindow.Show()
	w.fyneWindow.RequestFocus()
	w.input.FocusGained()
}

// Hide hides the window
func (w *Window) Hide() {
	w.fyneWindow.Hide()
}

// Close closes the window
func (w *Window) Close() {
	w.fyneWindow.Close()
}

// CenterOnScreen centers the window on screen
func (w *Window) CenterOnScreen() {
	w.fyneWindow.CenterOnScreen()
}

// SetOnClosed sets the callback to be called when the window is closed
func (w *Window) SetOnClosed(callback func()) {
	w.fyneWindow.SetOnClosed(callback)
}

// SetCloseIntercept sets the callback to be called when the window is about to close
func (w *Window) SetCloseIntercept(callback func()) {
	w.fyneWindow.SetCloseIntercept(callback)
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.fyneWindow
}
