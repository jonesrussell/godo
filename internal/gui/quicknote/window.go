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
	input  *debugEntry
}

// debugEntry extends Entry to add key press debugging
type debugEntry struct {
	widget.Entry
	logger logger.Logger
}

func newDebugEntry(log logger.Logger) *debugEntry {
	entry := &debugEntry{logger: log}
	entry.MultiLine = true
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *debugEntry) KeyDown(key *fyne.KeyEvent) {
	e.logger.Debug("Key pressed",
		"key", key.Name,
		"physical", key.Physical)
	e.Entry.KeyDown(key)
}

func (e *debugEntry) TypedShortcut(shortcut fyne.Shortcut) {
	e.logger.Debug("Shortcut received",
		"shortcut", fmt.Sprintf("%T", shortcut),
		"details", fmt.Sprintf("%+v", shortcut))
	e.Entry.TypedShortcut(shortcut)
}

func (e *debugEntry) TypedKey(key *fyne.KeyEvent) {
	e.logger.Debug("Key typed",
		"key", key.Name,
		"physical", key.Physical)
	e.Entry.TypedKey(key)
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

	// Create input field with debugging
	w.input = newDebugEntry(w.logger)
	w.input.SetPlaceHolder("Enter your quick note here...")

	// Create save button with shared save logic
	saveNote := func() {
		w.logger.Debug("Attempting to save note")
		text := w.input.Text
		if text == "" {
			w.logger.Debug("Note text is empty, skipping save")
			return
		}

		w.logger.Debug("Creating task from note", "text", text)
		// Create and save the task
		task := storage.Task{
			ID:        generateID(),
			Title:     text,
			Completed: false,
		}

		if err := w.store.Add(task); err != nil {
			w.logger.Error("Failed to save note", "error", err, "task", task)
			// TODO: Show error to user
			return
		}

		w.logger.Debug("Successfully saved note as task", "id", task.ID, "text", text)
		w.input.SetText("")
		w.logger.Debug("Cleared input field")
		w.Hide()
		w.logger.Debug("Hidden quick note window")
	}

	// Create save button
	saveBtn := widget.NewButton("Save", saveNote)

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

	// Add Shift+Enter shortcut for saving
	w.win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierShift,
	}, func(shortcut fyne.Shortcut) {
		w.logger.Debug("Shift+Enter shortcut triggered")
		if w.input == nil {
			w.logger.Error("Input field is nil when handling Shift+Enter")
			return
		}
		w.logger.Debug("Calling saveNote from Shift+Enter handler", "text", w.input.Text)
		saveNote()
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
