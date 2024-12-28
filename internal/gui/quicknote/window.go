package quicknote

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Interface defines the behavior of a quick note window
type Interface interface {
	Setup() error
	Hide()
	Show()
}

// Window implements the quick note window
type Window struct {
	store  storage.Store
	logger logger.Logger
	win    fyne.Window
	input  *widget.Entry
}

// New creates a new quick note window
func New(store storage.Store, logger logger.Logger) Interface {
	return newWindow(store, logger)
}

func newWindow(store storage.Store, logger logger.Logger) Interface {
	return &Window{
		store:  store,
		logger: logger,
	}
}

func (w *Window) Setup() error {
	w.logger.Debug("Setting up quick note window")

	// Create the window
	w.win = fyne.CurrentApp().NewWindow("Quick Note")
	w.win.Resize(fyne.NewSize(400, 300))
	w.win.CenterOnScreen()

	// Create input field
	w.input = widget.NewMultiLineEntry()
	w.input.SetPlaceHolder("Enter your quick note here...")

	// Create save button
	saveBtn := widget.NewButton("Save", func() {
		text := w.input.Text
		if text != "" {
			// Create and save the task
			task := storage.Task{
				ID:        generateID(),
				Title:     text,
				Completed: false,
			}

			if err := w.store.Add(task); err != nil {
				w.logger.Error("Failed to save note", "error", err)
				// TODO: Show error to user
				return
			}

			w.logger.Debug("Saved quick note as task", "id", task.ID)
			w.input.SetText("")
			w.Hide()
		}
	})

	// Create content container
	content := container.NewBorder(nil, saveBtn, nil, nil, w.input)
	w.win.SetContent(content)

	// Set window behavior
	w.win.SetCloseIntercept(func() {
		w.input.SetText("")
		w.Hide()
	})

	// Add ESC key handling
	w.win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyEscape,
	}, func(shortcut fyne.Shortcut) {
		w.logger.Debug("ESC key pressed, closing quick note window")
		w.input.SetText("")
		w.Hide()
	})

	w.logger.Debug("Quick note window setup complete")
	return nil
}

func (w *Window) Hide() {
	if w.win != nil {
		w.logger.Debug("Hiding quick note window")
		w.win.Hide()
	} else {
		w.logger.Error("Cannot hide quick note window - window is nil")
	}
}

func (w *Window) Show() {
	if w.win != nil {
		w.logger.Debug("Showing quick note window")
		w.win.Show()
		w.win.CenterOnScreen()
		if w.input != nil {
			w.win.Canvas().Focus(w.input)
		}
	} else {
		w.logger.Error("Cannot show quick note window - window is nil")
	}
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
