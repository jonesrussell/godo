//go:build !docker
// +build !docker

package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// window represents the quick note window
type window struct {
	app   fyne.App
	win   fyne.Window
	store storage.Store
	log   logger.Logger
	entry *customEntry
}

// newWindow creates a new quick note window
func newWindow(store storage.Store) Interface {
	w := &window{
		store: store,
	}
	return w
}

// Initialize sets up the window with the given app and logger
func (w *window) Initialize(app fyne.App, log logger.Logger) {
	w.app = app
	w.log = log
	w.win = app.NewWindow("Quick Note")
	w.entry = newCustomEntry(log)

	// Set close handler to hide window
	w.win.SetCloseIntercept(func() {
		w.Hide()
	})

	// Setup entry handlers
	w.entry.onCtrlEnter = func() {
		if w.entry.Text != "" {
			if err := w.store.SaveNote(w.entry.Text); err != nil {
				w.log.Error("Failed to save note", "error", err)
				return
			}
			w.log.Debug("Saved note", "content", w.entry.Text)
		}
		w.Hide()
	}
	w.entry.onEscape = func() {
		w.Hide()
	}

	// Setup window content
	content := container.NewVBox(
		widget.NewLabel("Enter your note (Ctrl+Enter to save, Esc to cancel):"),
		w.entry,
		container.NewHBox(
			widget.NewButton("Save", func() {
				w.entry.onCtrlEnter()
			}),
			widget.NewButton("Cancel", func() {
				w.entry.onEscape()
			}),
		),
	)

	w.win.SetContent(content)
	w.win.Resize(fyne.NewSize(400, 300))
	w.win.CenterOnScreen()
}

// Show displays the quick note window
func (w *window) Show() {
	w.entry.SetText("")
	w.win.Show()
	w.win.RequestFocus()
	w.win.Canvas().Focus(w.entry)
}

// Hide hides the quick note window
func (w *window) Hide() {
	w.entry.SetText("")
	w.win.Hide()
}
