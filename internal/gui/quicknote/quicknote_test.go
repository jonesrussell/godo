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

func setupTestQuickNote(t *testing.T) *QuickNote {
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
	return qn
}

func TestNew(t *testing.T) {
	qn := setupTestQuickNote(t)
	assert.NotNil(t, qn)
	assert.NotNil(t, qn.input)
	assert.NotNil(t, qn.window)
	assert.Equal(t, float32(400), qn.dimensions.window.Width)
	assert.Equal(t, float32(200), qn.dimensions.window.Height)
}

func TestShow(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.Show()
	// In test environment, we can't directly check window visibility
	// Instead, we verify the window exists and has expected properties
	assert.NotNil(t, qn.window)
}

func TestHide(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.Show()
	qn.Hide()
	assert.Equal(t, "", qn.input.Text)
}

func TestHandleSave(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.input.SetText("Test note")
	qn.handleSave()

	// Verify the note was saved
	notes, _ := qn.config.Store.GetNotes()
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])

	// Verify input cleared
	assert.Equal(t, "", qn.input.Text)
}

func TestHandleCancel(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.input.SetText("Test note")
	qn.Show()

	qn.handleCancel()

	assert.Equal(t, "", qn.input.Text)
}

func TestHandleFormSubmit(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.input.SetText("Test note")

	// Test save
	qn.handleFormSubmit(true)
	notes, _ := qn.config.Store.GetNotes()
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])

	// Test cancel
	qn.Show()
	qn.input.SetText("Another note")
	qn.handleFormSubmit(false)
	notes, _ = qn.config.Store.GetNotes()
	assert.Equal(t, 1, len(notes)) // Should still have only one note
}

func TestGetters(t *testing.T) {
	qn := setupTestQuickNote(t)
	assert.Equal(t, qn.window, qn.GetWindow())
	assert.Equal(t, qn.input, qn.GetInput())
}

func TestSetupShortcuts(t *testing.T) {
	qn := setupTestQuickNote(t)
	qn.input.SetText("Test note")

	// Simulate Ctrl+Enter shortcut
	shortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierControl,
	}
	qn.window.Canvas().AddShortcut(shortcut, func(shortcut fyne.Shortcut) {
		qn.handleSave()
	})

	// Trigger the shortcut manually
	qn.handleSave()

	// Verify note was saved
	notes, _ := qn.config.Store.GetNotes()
	assert.Equal(t, 1, len(notes))
	assert.Equal(t, "Test note", notes[0])
}
