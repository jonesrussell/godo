package quicknote

import (
	"fmt"

	"fyne.io/fyne/v2"
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
func New(store storage.Store, l logger.Logger) Interface {
	return &Window{
		store:  store,
		logger: l,
	}
}

// Setup initializes the window
func (w *Window) Setup() error {
	w.logger.Debug("Setting up quick note window")

	// Create window
	w.win = fyne.CurrentApp().NewWindow("Quick Note")
	if w.win == nil {
		w.logger.Error("Failed to create window")
		return fmt.Errorf("failed to create window")
	}

	// Create input field
	w.input = newDebugEntry(w.logger)

	// Set window content
	w.win.SetContent(w.input)

	// Set window properties
	w.win.Resize(fyne.NewSize(400, 200))
	w.win.CenterOnScreen()

	// Add keyboard shortcuts
	w.win.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyEscape,
	}, func(_ fyne.Shortcut) {
		w.input.SetText("")
		w.Hide()
	})

	return nil
}

// Hide hides the window
func (w *Window) Hide() {
	if w.win != nil {
		w.win.Hide()
	} else {
		w.logger.Error("Cannot hide window: window not initialized")
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
