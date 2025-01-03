package mainwindow

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// MainWindow represents the main application window
type MainWindow struct {
	window fyne.Window
	store  types.Store
	list   *widget.List
	notes  []types.Note
}

// NewMainWindow creates a new main window
func NewMainWindow(window fyne.Window, store types.Store) *MainWindow {
	w := &MainWindow{
		window: window,
		store:  store,
	}
	w.createUI()
	return w
}

// Show displays the main window
func (w *MainWindow) Show() {
	w.window.Show()
}

// Hide hides the main window
func (w *MainWindow) Hide() {
	w.window.Hide()
}

// Close closes the main window
func (w *MainWindow) Close() {
	w.window.Close()
}

func (w *MainWindow) createUI() {
	// Create the note list
	w.list = widget.NewList(
		func() int {
			return len(w.notes)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			note := w.notes[id]
			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)
			deleteBtn := box.Objects[2].(*widget.Button)

			check.Checked = note.Completed
			check.OnChanged = func(checked bool) {
				w.toggleNote(id, checked)
			}

			label.SetText(note.Content)
			if note.Completed {
				label.TextStyle = fyne.TextStyle{Monospace: true}
			} else {
				label.TextStyle = fyne.TextStyle{}
			}

			deleteBtn.OnTapped = func() {
				w.deleteNote(id)
			}
		},
	)

	// Create the add note button
	addBtn := widget.NewButtonWithIcon("Add Note", theme.ContentAddIcon(), w.showAddNoteDialog)

	// Create the main container
	content := container.NewBorder(nil, addBtn, nil, nil, w.list)
	w.window.SetContent(content)

	// Load initial notes
	w.refreshNotes()
}

func (w *MainWindow) refreshNotes() {
	notes, err := w.store.List(context.Background())
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load notes: %w", err), w.window)
		return
	}
	w.notes = notes
	w.list.Refresh()
}

func (w *MainWindow) showAddNoteDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter note content")

	dialog.ShowForm("Add Note", "Add", "Cancel",
		[]*widget.FormItem{
			{Text: "Content", Widget: entry},
		},
		func(confirm bool) {
			if !confirm || entry.Text == "" {
				return
			}

			note := types.Note{
				ID:        uuid.New().String(),
				Content:   entry.Text,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			if err := w.store.Add(context.Background(), note); err != nil {
				dialog.ShowError(fmt.Errorf("failed to add note: %w", err), w.window)
				return
			}

			w.refreshNotes()
		},
		w.window,
	)
}

func (w *MainWindow) toggleNote(id widget.ListItemID, checked bool) {
	note := w.notes[id]
	note.Completed = checked
	note.UpdatedAt = time.Now().Unix()

	if err := w.store.Update(context.Background(), note); err != nil {
		dialog.ShowError(fmt.Errorf("failed to update note: %w", err), w.window)
		return
	}

	w.refreshNotes()
}

func (w *MainWindow) deleteNote(id widget.ListItemID) {
	note := w.notes[id]

	dialog.ShowConfirm("Delete Note",
		fmt.Sprintf("Are you sure you want to delete note: %s?", note.Content),
		func(confirm bool) {
			if !confirm {
				return
			}

			if err := w.store.Delete(context.Background(), note.ID); err != nil {
				dialog.ShowError(fmt.Errorf("failed to delete note: %w", err), w.window)
				return
			}

			w.refreshNotes()
		},
		w.window,
	)
}
