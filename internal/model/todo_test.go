package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *Todo
	}{
		{
			name:    "creates todo with content",
			content: "Test todo",
			want: &Todo{
				Content: "Test todo",
				Done:    false,
			},
		},
		{
			name:    "creates todo with empty content",
			content: "",
			want: &Todo{
				Content: "",
				Done:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTodo(tt.content)

			// Verify UUID format
			_, err := uuid.Parse(got.ID)
			assert.NoError(t, err)

			// Verify content and done status
			assert.Equal(t, tt.want.Content, got.Content)
			assert.Equal(t, tt.want.Done, got.Done)

			// Verify timestamps
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)
			assert.Equal(t, got.CreatedAt, got.UpdatedAt)

			// Verify timestamps are recent
			now := time.Now()
			assert.True(t, got.CreatedAt.Before(now) || got.CreatedAt.Equal(now))
			assert.True(t, got.UpdatedAt.Before(now) || got.UpdatedAt.Equal(now))
		})
	}
}

func TestTodo_ToggleDone(t *testing.T) {
	tests := []struct {
		name     string
		todo     *Todo
		wantDone bool
	}{
		{
			name: "toggles undone to done",
			todo: &Todo{
				ID:        uuid.New().String(),
				Content:   "Test todo",
				Done:      false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantDone: true,
		},
		{
			name: "toggles done to undone",
			todo: &Todo{
				ID:        uuid.New().String(),
				Content:   "Test todo",
				Done:      true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantDone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			tt.todo.ToggleDone()

			assert.Equal(t, tt.wantDone, tt.todo.Done)
			assert.True(t, tt.todo.UpdatedAt.After(originalUpdatedAt))
		})
	}
}

func TestTodo_UpdateContent(t *testing.T) {
	tests := []struct {
		name        string
		todo        *Todo
		newContent  string
		wantContent string
	}{
		{
			name: "updates with new content",
			todo: &Todo{
				ID:        uuid.New().String(),
				Content:   "Original content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			newContent:  "Updated content",
			wantContent: "Updated content",
		},
		{
			name: "updates with empty content",
			todo: &Todo{
				ID:        uuid.New().String(),
				Content:   "Original content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			newContent:  "",
			wantContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			tt.todo.UpdateContent(tt.newContent)

			assert.Equal(t, tt.wantContent, tt.todo.Content)
			assert.True(t, tt.todo.UpdatedAt.After(originalUpdatedAt))
		})
	}
}
