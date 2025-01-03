package validation

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestValidateNote(t *testing.T) {
	t.Run("valid note", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateNote(note)
		assert.NoError(t, err)
	})

	t.Run("empty ID", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateNote(note)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID cannot be empty")
	})

	t.Run("empty content", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := ValidateNote(note)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "content cannot be empty")
	})

	t.Run("zero created_at", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: 0,
			UpdatedAt: now,
		}

		err := ValidateNote(note)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "created_at cannot be zero")
	})

	t.Run("zero updated_at", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: 0,
		}

		err := ValidateNote(note)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "updated_at cannot be zero")
	})

	t.Run("updated_at before created_at", func(t *testing.T) {
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now - 1,
		}

		err := ValidateNote(note)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "updated_at cannot be before created_at")
	})
}

func TestValidateID(t *testing.T) {
	t.Run("valid ID", func(t *testing.T) {
		err := ValidateID("test-1")
		assert.NoError(t, err)
	})

	t.Run("empty ID", func(t *testing.T) {
		err := ValidateID("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidID, err)
	})

	t.Run("too long ID", func(t *testing.T) {
		longID := make([]byte, MaxIDLength+1)
		for i := range longID {
			longID[i] = 'a'
		}
		err := ValidateID(string(longID))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidID, err)
	})
}

func TestValidateContent(t *testing.T) {
	t.Run("valid content", func(t *testing.T) {
		err := ValidateContent("Test Note")
		assert.NoError(t, err)
	})

	t.Run("empty content", func(t *testing.T) {
		err := ValidateContent("")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidContent, err)
	})

	t.Run("too long content", func(t *testing.T) {
		longContent := make([]byte, MaxContentLength+1)
		for i := range longContent {
			longContent[i] = 'a'
		}
		err := ValidateContent(string(longContent))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidContent, err)
	})
}
