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

// Window represents the quick note window
type Window struct {
	app   fyne.App
	win   fyne.Window
	store storage.Store
	log   logger.Logger
	entry *customEntry
}

// New creates a new quick note window
func New(app fyne.App, store storage.Store, log logger.Logger) *Window {
	win := app.NewWindow("Quick Note")

	w := &Window{
		app:   app,
		win:   win,
		store: store,
		log:   log,
		entry: newCustomEntry(log),
	}

	// Set close handler to hide window
	win.SetCloseIntercept(func() {
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

	win.SetContent(content)
	win.Resize(fyne.NewSize(300, 100))
	win.CenterOnScreen()

	return w
}

// Show displays the quick note window
func (w *Window) Show() {
	w.entry.SetText("")
	w.win.Show()
	w.win.RequestFocus()
	w.win.Canvas().Focus(w.entry)
}

// Hide hides the quick note window
func (w *Window) Hide() {
	w.entry.SetText("")
	w.win.Hide()
}
