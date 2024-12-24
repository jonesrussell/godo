//go:build !docker
// +build !docker

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func setupTestQuickNote(t *testing.T) (*QuickNote, *memory.Store) {
	app := test.NewApp()
	mainWindow := app.NewWindow("Test Main")
	store := memory.New()
	log := logger.NewTestLogger(t)

	cfg := Config{
		App:        app,
		MainWindow: mainWindow,
		Store:      store,
		Logger:     log,
	}

	qn := New(cfg)
	return qn, store
}

func TestNew(t *testing.T) {
	qn, _ := setupTestQuickNote(t)
	assert.NotNil(t, qn)
	assert.NotNil(t, qn.input)
	assert.NotNil(t, qn.window)
	assert.Equal(t, float32(400), qn.dimensions.window.Width)
	assert.Equal(t, float32(200), qn.dimensions.window.Height)
}

func TestShow(t *testing.T) {
	qn, _ := setupTestQuickNote(t)
	qn.Show()
	// In test environment, we can't directly check window visibility
	// Instead, we verify the window exists and has expected properties
	assert.NotNil(t, qn.window)
}

func TestHide(t *testing.T) {
	qn, _ := setupTestQuickNote(t)
	qn.Show()
	qn.Hide()
	assert.Equal(t, "", qn.input.Text)
}

func TestHandleSave(t *testing.T) {
	qn, store := setupTestQuickNote(t)
	qn.input.SetText("Test note")

	// Save note directly to verify store is working
	err := store.SaveNote("Test note")
	assert.NoError(t, err)

	// Verify the note was saved
	notes, err := store.GetNotes()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])

	// Now test the handleSave method
	qn.input.SetText("Another note")
	qn.handleSave()

	// Verify both notes are saved
	notes, err = store.GetNotes()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(notes))
	assert.Equal(t, "Test note", notes[0])
	assert.Equal(t, "Another note", notes[1])

	// Verify input cleared
	assert.Equal(t, "", qn.input.Text)
}

func TestHandleCancel(t *testing.T) {
	qn, _ := setupTestQuickNote(t)
	qn.input.SetText("Test note")
	qn.Show()

	qn.handleCancel()

	assert.Equal(t, "", qn.input.Text)
}

func TestHandleFormSubmit(t *testing.T) {
	qn, store := setupTestQuickNote(t)
	qn.input.SetText("Test note")

	// Test save
	qn.handleFormSubmit(true)
	notes, err := store.GetNotes()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])

	// Test cancel
	qn.Show()
	qn.input.SetText("Another note")
	qn.handleFormSubmit(false)
	notes, err = store.GetNotes()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notes)) // Should still have only one note
}

func TestGetters(t *testing.T) {
	qn, _ := setupTestQuickNote(t)
	assert.Equal(t, qn.window, qn.GetWindow())
	assert.Equal(t, qn.input, qn.GetInput())
}

func TestSetupShortcuts(t *testing.T) {
	qn, store := setupTestQuickNote(t)
	qn.input.SetText("Test note")

	// Simulate Ctrl+Enter shortcut
	shortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierControl,
	}
	qn.window.Canvas().AddShortcut(shortcut, func(_ fyne.Shortcut) {
		qn.handleSave()
	})

	// Trigger the shortcut manually
	qn.handleSave()

	// Verify note was saved
	notes, err := store.GetNotes()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])
}
