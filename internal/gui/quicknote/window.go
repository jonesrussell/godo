// Package quicknote implements the quick note window functionality
package quicknote

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

// Window represents the quick note window
type Window struct {
	fyneWindow fyne.Window
	store      storage.Store
	logger     logger.Logger
	config     config.WindowConfig
	input      *widget.Entry
	saveBtn    *widget.Button
}

// New creates a new quick note window instance
func New(app fyne.App, store storage.Store, logger logger.Logger, cfg config.WindowConfig) *Window {
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
		if w.input.Text == "" {
			return
		}

		now := time.Now().Unix()
		task := storage.Task{
			ID:        fmt.Sprintf("%d", now), // Simple ID generation
			Title:     w.input.Text,
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := w.store.Add(context.Background(), task); err != nil {
			w.logger.Error("Failed to save quick note", "error", err)
			return
		}

		w.input.SetText("")
		w.Hide()
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
