//go:build !docker

package app

import (
	"context"
	"fmt"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStore implements types.Store for testing
type mockStore struct {
	notes map[string]types.Note
	err   error
}

func newMockStore() *mockStore {
	return &mockStore{
		notes: make(map[string]types.Note),
	}
}

func (s *mockStore) Add(_ context.Context, note types.Note) error {
	if s.err != nil {
		return s.err
	}
	s.notes[note.ID] = note
	return nil
}

func (s *mockStore) Get(_ context.Context, id string) (types.Note, error) {
	if s.err != nil {
		return types.Note{}, s.err
	}
	note, ok := s.notes[id]
	if !ok {
		return types.Note{}, fmt.Errorf("note not found")
	}
	return note, nil
}

func (s *mockStore) List(_ context.Context) ([]types.Note, error) {
	if s.err != nil {
		return nil, s.err
	}
	notes := make([]types.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return notes, nil
}

func (s *mockStore) Update(_ context.Context, note types.Note) error {
	if s.err != nil {
		return s.err
	}
	if _, ok := s.notes[note.ID]; !ok {
		return fmt.Errorf("note not found")
	}
	s.notes[note.ID] = note
	return nil
}

func (s *mockStore) Delete(_ context.Context, id string) error {
	if s.err != nil {
		return s.err
	}
	if _, ok := s.notes[id]; !ok {
		return fmt.Errorf("note not found")
	}
	delete(s.notes, id)
	return nil
}

func (s *mockStore) BeginTx(_ context.Context) (types.Transaction, error) {
	return nil, fmt.Errorf("transactions not supported")
}

func (s *mockStore) Close() error {
	return s.err
}

// mockWindow implements gui.MainWindowManager for testing
type mockWindow struct {
	content      fyne.CanvasObject
	showCalled   bool
	hideCalled   bool
	sizeCalled   bool
	centerCalled bool
}

func newMockWindow() *mockWindow {
	return &mockWindow{}
}

func (w *mockWindow) Show() {
	w.showCalled = true
}

func (w *mockWindow) Hide() {
	w.hideCalled = true
}

func (w *mockWindow) CenterOnScreen() {
	w.centerCalled = true
}

func (w *mockWindow) SetContent(content fyne.CanvasObject) {
	w.content = content
}

func (w *mockWindow) Resize(size fyne.Size) {
	w.sizeCalled = true
}

func (w *mockWindow) GetWindow() fyne.Window {
	return test.NewWindow(nil)
}

func TestAppOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("New", func(t *testing.T) {
		store := newMockStore()
		window := newMockWindow()

		app, err := New(Params{
			Store:  store,
			Window: window,
		})
		require.NoError(t, err)
		assert.NotNil(t, app)
	})

	t.Run("Start", func(t *testing.T) {
		store := newMockStore()
		window := newMockWindow()

		app, err := New(Params{
			Store:  store,
			Window: window,
		})
		require.NoError(t, err)

		err = app.Start()
		require.NoError(t, err)
		assert.True(t, window.showCalled)
		assert.NotNil(t, window.content)
	})

	t.Run("AddNote", func(t *testing.T) {
		store := newMockStore()
		window := newMockWindow()

		app, err := New(Params{
			Store:  store,
			Window: window,
		})
		require.NoError(t, err)
		require.NoError(t, app.Start())

		note := types.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err = app.AddNote(ctx, note)
		require.NoError(t, err)

		// Verify note was added to store
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, note.Content, notes[0].Content)
	})

	t.Run("UpdateNote", func(t *testing.T) {
		store := newMockStore()
		window := newMockWindow()

		app, err := New(Params{
			Store:  store,
			Window: window,
		})
		require.NoError(t, err)
		require.NoError(t, app.Start())

		note := types.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		// Add note first
		err = app.AddNote(ctx, note)
		require.NoError(t, err)

		// Update note
		note.Content = "Updated Note"
		note.Completed = true
		err = app.UpdateNote(ctx, note)
		require.NoError(t, err)

		// Verify note was updated
		updated, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.Content, updated.Content)
		assert.Equal(t, note.Completed, updated.Completed)
	})

	t.Run("DeleteNote", func(t *testing.T) {
		store := newMockStore()
		window := newMockWindow()

		app, err := New(Params{
			Store:  store,
			Window: window,
		})
		require.NoError(t, err)
		require.NoError(t, app.Start())

		note := types.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		// Add note first
		err = app.AddNote(ctx, note)
		require.NoError(t, err)

		// Delete note
		err = app.DeleteNote(ctx, note.ID)
		require.NoError(t, err)

		// Verify note was deleted
		_, err = store.Get(ctx, note.ID)
		assert.Error(t, err)
	})
}
