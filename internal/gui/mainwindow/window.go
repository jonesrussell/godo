// Package mainwindow provides the main application window
package mainwindow

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// Window represents the main application window
type Window struct {
	fyneWindow fyne.Window
	store      types.Store
	logger     logger.Logger
	list       *widget.List
	notes      map[string]*note.Note
}

// New creates a new main window
func New(win fyne.Window, store types.Store, log logger.Logger) *Window {
	w := &Window{
		fyneWindow: win,
		store:      store,
		logger:     log,
		notes:      make(map[string]*note.Note),
	}

	w.setupUI()
	w.loadNotes()

	return w
}

// setupUI initializes the window UI
func (w *Window) setupUI() {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter a note...")

	w.list = widget.NewList(
		func() int { return len(w.notes) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			note := w.getNoteByIndex(id)
			if note == nil {
				return
			}

			box := obj.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			check.Checked = note.Completed
			check.OnChanged = func(checked bool) {
				note.Completed = checked
				note.UpdatedAt = time.Now()
				w.updateNote(note)
			}

			label.SetText(note.Content)
		},
	)

	input.OnSubmitted = func(text string) {
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
			w.logger.Error("Failed to add note", "error", err)
			return
		}

		w.notes[note.ID] = note
		w.list.Refresh()
		input.SetText("")
	}

	content := container.NewBorder(input, nil, nil, nil, w.list)
	w.fyneWindow.SetContent(content)
	w.fyneWindow.Resize(fyne.NewSize(400, 300))
}

// loadNotes loads notes from storage
func (w *Window) loadNotes() {
	notes, err := w.store.List(context.Background())
	if err != nil {
		w.logger.Error("Failed to load notes", "error", err)
		return
	}

	for _, note := range notes {
		w.notes[note.ID] = note
	}
	w.list.Refresh()
}

// updateNote updates a note in storage
func (w *Window) updateNote(note *note.Note) {
	if err := w.store.Update(context.Background(), note); err != nil {
		w.logger.Error("Failed to update note", "error", err)
		return
	}
	w.notes[note.ID] = note
	w.list.Refresh()
}

// getNoteByIndex returns a note by its list index
func (w *Window) getNoteByIndex(index int) *note.Note {
	i := 0
	for _, note := range w.notes {
		if i == index {
			return note
		}
		i++
	}
	return nil
}

// CenterOnScreen centers the window on screen
func (w *Window) CenterOnScreen() {
	w.fyneWindow.CenterOnScreen()
}

// Show shows the window
func (w *Window) Show() {
	w.fyneWindow.Show()
}

// Hide hides the window
func (w *Window) Hide() {
	w.fyneWindow.Hide()
}

// GetWindow returns the underlying fyne.Window
func (w *Window) GetWindow() fyne.Window {
	return w.fyneWindow
}

// Resize resizes the window
func (w *Window) Resize(size fyne.Size) {
	w.fyneWindow.Resize(size)
}

// SetContent sets the window content
func (w *Window) SetContent(content fyne.CanvasObject) {
	w.fyneWindow.SetContent(content)
}
