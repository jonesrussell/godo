//go:build integration
// +build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteHandler(t *testing.T) {
	store := storage.NewMockStore()
	handler := NewNoteHandler(store)

	t.Run("create note", func(t *testing.T) {
		tests := []struct {
			name    string
			note    storage.Note
			wantErr bool
		}{
			{
				name: "valid note",
				note: storage.Note{
					Content:   "Test Note",
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				wantErr: false,
			},
			{
				name:    "invalid note",
				note:    storage.Note{},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.note)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
				w := httptest.NewRecorder()

				handler.CreateNote(w, req)

				if tt.wantErr {
					assert.Equal(t, http.StatusBadRequest, w.Code)
					return
				}

				assert.Equal(t, http.StatusCreated, w.Code)

				var response storage.Note
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.note.Content, response.Content)
				assert.Equal(t, tt.note.CreatedAt, response.CreatedAt)
				assert.Equal(t, tt.note.UpdatedAt, response.UpdatedAt)
			})
		}
	})

	t.Run("list notes", func(t *testing.T) {
		// Clear any existing notes
		store = storage.NewMockStore()

		// Add some test notes
		notes := []storage.Note{
			{
				ID:        "test-1",
				Content:   "Note 1",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			{
				ID:        "test-2",
				Content:   "Note 2",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
		}

		for _, note := range notes {
			err := store.Add(context.Background(), note)
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()

		handler.ListNotes(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []storage.Note
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response, len(notes))

		for i, note := range notes {
			assert.Equal(t, note.Content, response[i].Content)
			assert.Equal(t, note.CreatedAt, response[i].CreatedAt)
			assert.Equal(t, note.UpdatedAt, response[i].UpdatedAt)
		}
	})

	t.Run("update note", func(t *testing.T) {
		tests := []struct {
			name    string
			noteID  string
			update  storage.Note
			wantErr bool
		}{
			{
				name:   "existing note",
				noteID: "test-1",
				update: storage.Note{
					ID:        "test-1",
					Content:   "Updated Note",
					UpdatedAt: time.Now().Unix(),
				},
				wantErr: false,
			},
			{
				name:   "nonexistent note",
				noteID: "nonexistent",
				update: storage.Note{
					ID:        "nonexistent",
					Content:   "Updated Note",
					UpdatedAt: time.Now().Unix(),
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.noteID == "test-1" {
					existingNote := storage.Note{
						ID:        "test-1",
						Content:   "Original Note",
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
					}
					err := store.Add(context.Background(), existingNote)
					require.NoError(t, err)
				}

				body, err := json.Marshal(tt.update)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPut, "/notes/"+tt.noteID, bytes.NewReader(body))
				w := httptest.NewRecorder()

				handler.UpdateNote(w, req)

				if tt.wantErr {
					assert.Equal(t, http.StatusNotFound, w.Code)
					return
				}

				assert.Equal(t, http.StatusOK, w.Code)

				var response storage.Note
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.update.Content, response.Content)
				assert.Equal(t, tt.update.UpdatedAt, response.UpdatedAt)
			})
		}
	})

	t.Run("delete note", func(t *testing.T) {
		tests := []struct {
			name    string
			noteID  string
			wantErr bool
		}{
			{
				name:    "existing note",
				noteID:  "test-1",
				wantErr: false,
			},
			{
				name:    "nonexistent note",
				noteID:  "nonexistent",
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.noteID == "test-1" {
					note := storage.Note{
						ID:        "test-1",
						Content:   "Test Note",
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
					}
					err := store.Add(context.Background(), note)
					require.NoError(t, err)
				}

				req := httptest.NewRequest(http.MethodDelete, "/notes/"+tt.noteID, nil)
				w := httptest.NewRecorder()

				handler.DeleteNote(w, req)

				if tt.wantErr {
					assert.Equal(t, http.StatusNotFound, w.Code)
					return
				}

				assert.Equal(t, http.StatusNoContent, w.Code)
			})
		}
	})
}
