package gui

import (
	"context"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockWindow(t *testing.T) {
	testWindow := test.NewWindow(nil)

	tests := []struct {
		name  string
		title string
	}{
		{
			name:  "creates window with valid title",
			title: "Test Window",
		},
		{
			name:  "creates window with empty title",
			title: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock window
			win := &MockMainWindow{
				Window: testWindow,
			}

			assert.NotNil(t, win)
			assert.NotNil(t, win.GetWindow())

			// Test window properties
			assert.False(t, win.ShowCalled)
			assert.False(t, win.HideCalled)
			assert.False(t, win.ResizeCalled)
			assert.False(t, win.CenterCalled)
			assert.Nil(t, win.ContentSet)

			// Test window operations
			win.Show()
			assert.True(t, win.ShowCalled)

			win.Hide()
			assert.True(t, win.HideCalled)

			win.Resize(fyne.NewSize(800, 600))
			assert.True(t, win.ResizeCalled)

			win.CenterOnScreen()
			assert.True(t, win.CenterCalled)

			content := container.NewVBox()
			win.SetContent(content)
			assert.Equal(t, content, win.ContentSet)
		})
	}
}

func TestMockWindowWithNotes(t *testing.T) {
	// Create test dependencies
	store := storage.NewMockStore()
	testWindow := test.NewWindow(nil)
	ctx := context.Background()

	// Create mock window
	win := &MockMainWindow{
		Window: testWindow,
	}

	// Add a test note
	note := storage.Note{
		ID:        "test-note",
		Content:   "Test Note",
		Completed: false,
	}
	err := store.Add(ctx, note)
	require.NoError(t, err)

	// Test note operations
	content := container.NewVBox()
	win.SetContent(content)
	assert.Equal(t, content, win.ContentSet)

	// Test note completion
	note.Completed = true
	err = store.Update(ctx, note)
	require.NoError(t, err)

	// Test note deletion
	err = store.Delete(ctx, note.ID)
	require.NoError(t, err)
}
