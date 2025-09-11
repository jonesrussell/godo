// Package mainwindow provides the main application window
package mainwindow

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

//go:generate mockgen -destination=../../../test/mocks/mock_mainwindow.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/gui/mainwindow Interface

// Interface defines the main window functionality
type Interface interface {
	Show()
	Hide()
	SetContent(content fyne.CanvasObject)
	Resize(size fyne.Size)
	CenterOnScreen()
	GetWindow() fyne.Window
	Refresh()
}

// Window represents the main application window
type Window struct {
	app    fyne.App
	window fyne.Window
	store  storage.NoteStore
	log    logger.Logger
	notes  []model.Note
	cfg    config.WindowConfig

	// UI components
	noteList    *widget.List
	addButton   *widget.Button
	refreshBtn  *widget.Button
	searchEntry *widget.Entry
	toolbar     *fyne.Container
	statusBar   *widget.Label
}

// New creates a new main window
func New(app fyne.App, store storage.NoteStore, log logger.Logger, cfg config.WindowConfig) *Window {
	w := &Window{
		app:    app,
		store:  store,
		log:    log,
		cfg:    cfg,
		notes:  make([]model.Note, 0),
		window: app.NewWindow("Godo - Note Manager"),
	}

	w.setupUI()
	w.loadNotes()
	return w
}

// Show displays the main window
func (w *Window) Show() {
	fyne.Do(func() {
		w.window.Show()
	})
}

// Hide hides the main window
func (w *Window) Hide() {
	fyne.Do(func() {
		w.window.Hide()
	})
}

// GetWindow returns the underlying Fyne window
func (w *Window) GetWindow() fyne.Window {
	return w.window
}

// SetContent sets the content of the main window
func (w *Window) SetContent(content fyne.CanvasObject) {
	fyne.Do(func() {
		w.window.SetContent(content)
	})
}

// setupUI initializes the user interface
func (w *Window) setupUI() {
	w.createNoteList()
	w.createToolbar()
	w.createStatusBar()
	w.createMainLayout()
}

// createNoteList creates the note list widget
func (w *Window) createNoteList() {
	w.noteList = widget.NewList(
		func() int { return len(w.notes) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel("Note content"),
				layout.NewSpacer(),
				widget.NewButton("Edit", nil),
				widget.NewButton("Delete", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			w.updateNoteListItem(id, obj)
		},
	)
}

// updateNoteListItem updates a note list item
func (w *Window) updateNoteListItem(id widget.ListItemID, obj fyne.CanvasObject) {
	box, ok := obj.(*fyne.Container)
	if !ok {
		w.log.Error("Failed to cast object to container")
		return
	}

	if id >= len(w.notes) {
		return
	}
	note := w.notes[id]

	// Update check box
	if check, okCheck := box.Objects[0].(*widget.Check); okCheck {
		check.Checked = note.Done
		check.OnChanged = func(checked bool) {
			w.toggleNote(id, checked)
		}
	}

	// Update label
	if label, okLabel := box.Objects[1].(*widget.Label); okLabel {
		label.SetText(note.Content)
	}

	// Update edit button
	if editBtn, okEdit := box.Objects[3].(*widget.Button); okEdit {
		editBtn.OnTapped = func() {
			w.editNote(id)
		}
	}

	// Update delete button
	if deleteBtn, okDelete := box.Objects[4].(*widget.Button); okDelete {
		deleteBtn.OnTapped = func() {
			w.deleteNote(id)
		}
	}
}

// createToolbar creates the toolbar with buttons and search
func (w *Window) createToolbar() {
	// Create add button
	w.addButton = widget.NewButtonWithIcon("Add Note", theme.ContentAddIcon(), w.addNote)

	// Create refresh button
	w.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), w.loadNotes)

	// Create search entry
	w.searchEntry = widget.NewEntry()
	w.searchEntry.SetPlaceHolder("Search notes...")
	w.searchEntry.OnChanged = w.filterNotes

	// Create toolbar
	w.toolbar = container.NewHBox(
		w.addButton,
		w.refreshBtn,
		layout.NewSpacer(),
		w.searchEntry,
	)
}

// createStatusBar creates the status bar
func (w *Window) createStatusBar() {
	w.statusBar = widget.NewLabel("")
	w.statusBar.Hide()
}

// createMainLayout creates the main window layout
func (w *Window) createMainLayout() {
	fyne.Do(func() {
		// Create main container
		content := container.NewBorder(
			w.toolbar,
			w.statusBar,
			nil,
			nil,
			w.noteList,
		)

		w.window.SetContent(content)
		w.window.Resize(fyne.NewSize(600, 400))
	})
}

// loadNotes loads all notes from storage
func (w *Window) loadNotes() {
	ctx := context.Background()
	notes, err := w.store.List(ctx)
	if err != nil {
		w.log.Error("Failed to load notes", "error", err)
		w.showStatus("Failed to load notes", true)
		return
	}

	w.notes = notes
	w.noteList.Refresh()
	w.showStatus(fmt.Sprintf("Loaded %d notes", len(notes)), false)
	w.log.Info("Notes loaded", "count", len(notes))
}

// addNote adds a new note with enhanced dialog
func (w *Window) addNote() {
	content := widget.NewEntry()
	content.SetPlaceHolder("Enter note content...")

	form := dialog.NewForm(
		"Add Note",
		"Add",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Note", content),
		},
		func(confirm bool) {
			if !confirm || content.Text == "" {
				return
			}

			note := model.Note{
				ID:        uuid.New().String(),
				Content:   content.Text,
				Done:      false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			ctx := context.Background()
			if err := w.store.Add(ctx, &note); err != nil {
				w.log.Error("Failed to add note", "error", err)
				w.showStatus("Failed to add note", true)
				return
			}

			w.loadNotes()
			w.showStatus("Note added successfully", false)
		},
		w.window,
	)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
}

// toggleNote toggles the done status of a note
func (w *Window) toggleNote(id widget.ListItemID, done bool) {
	if id >= len(w.notes) {
		return
	}

	note := w.notes[id]
	note.Done = done
	note.UpdatedAt = time.Now()

	ctx := context.Background()
	if err := w.store.Update(ctx, &note); err != nil {
		w.log.Error("Failed to update note", "note_id", note.ID, "error", err)
		w.showStatus("Failed to update note", true)
		return
	}

	w.notes[id] = note
	w.noteList.Refresh()
	w.showStatus("Note updated", false)
}

// editNote opens a dialog to edit a note
func (w *Window) editNote(id widget.ListItemID) {
	if id >= len(w.notes) {
		return
	}

	note := w.notes[id]
	content := widget.NewEntry()
	content.SetText(note.Content)

	form := dialog.NewForm(
		"Edit Note",
		"Save",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Note", content),
		},
		func(confirm bool) {
			if !confirm || content.Text == "" {
				return
			}

			note.Content = content.Text
			note.UpdatedAt = time.Now()

			ctx := context.Background()
			if err := w.store.Update(ctx, &note); err != nil {
				w.log.Error("Failed to update note", "note_id", note.ID, "error", err)
				w.showStatus("Failed to update note", true)
				return
			}

			w.notes[id] = note
			w.noteList.Refresh()
			w.showStatus("Note updated", false)
		},
		w.window,
	)

	form.Resize(fyne.NewSize(400, 200))
	form.Show()
}

// deleteNote deletes a note with confirmation
func (w *Window) deleteNote(id widget.ListItemID) {
	if id >= len(w.notes) {
		return
	}

	note := w.notes[id]

	confirm := dialog.NewConfirm(
		"Delete Note",
		fmt.Sprintf("Are you sure you want to delete '%s'?", note.Content),
		func(confirm bool) {
			if !confirm {
				return
			}

			ctx := context.Background()
			if err := w.store.Delete(ctx, note.ID); err != nil {
				w.log.Error("Failed to delete note", "note_id", note.ID, "error", err)
				w.showStatus("Failed to delete note", true)
				return
			}

			w.loadNotes()
			w.showStatus("Note deleted", false)
		},
		w.window,
	)

	confirm.Show()
}

// filterNotes filters the note list based on search text
func (w *Window) filterNotes(searchText string) {
	// For now, just reload all notes
	// In a real implementation, you might want to implement client-side filtering
	w.loadNotes()
}

// showStatus shows a status message
func (w *Window) showStatus(message string, isError bool) {
	fyne.Do(func() {
		w.statusBar.SetText(message)
		w.statusBar.Show()

		if isError {
			w.statusBar.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			w.statusBar.TextStyle = fyne.TextStyle{}
		}
	})

	// Auto-hide after 3 seconds
	go func() {
		time.Sleep(3 * time.Second)
		fyne.Do(func() {
			w.statusBar.Hide()
		})
	}()
}

// CenterOnScreen centers the main window on the screen
func (w *Window) CenterOnScreen() {
	fyne.Do(func() {
		w.window.CenterOnScreen()
	})
}

// Refresh refreshes the main window content
func (w *Window) Refresh() {
	fyne.Do(func() {
		w.noteList.Refresh()
	})
}

// Resize resizes the main window
func (w *Window) Resize(size fyne.Size) {
	fyne.Do(func() {
		w.window.Resize(size)
	})
}
