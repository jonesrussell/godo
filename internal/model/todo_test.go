package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  *Todo
	}{
		{
			name:  "creates todo with title",
			title: "Test todo",
			want: &Todo{
				Title:     "Test todo",
				Completed: false,
			},
		},
		{
			name:  "creates todo with empty title",
			title: "",
			want: &Todo{
				Title:     "",
				Completed: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTodo(tt.title)

			// Verify UUID format
			_, err := uuid.Parse(got.ID)
			assert.NoError(t, err)

			// Verify title and completed status
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Completed, got.Completed)

			// Verify timestamps
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)
			assert.Equal(t, got.CreatedAt, got.UpdatedAt)

			// Verify timestamps are recent
			now := time.Now().Unix()
			assert.True(t, got.CreatedAt <= now)
			assert.True(t, got.UpdatedAt <= now)
		})
	}
}

func TestTodo_ToggleCompleted(t *testing.T) {
	tests := []struct {
		name          string
		todo          *Todo
		wantCompleted bool
	}{
		{
			name: "toggles uncompleted to completed",
			todo: &Todo{
				ID:        uuid.New().String(),
				Title:     "Test todo",
				Completed: false,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			wantCompleted: true,
		},
		{
			name: "toggles completed to uncompleted",
			todo: &Todo{
				ID:        uuid.New().String(),
				Title:     "Test todo",
				Completed: true,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			wantCompleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			tt.todo.ToggleCompleted()

			assert.Equal(t, tt.wantCompleted, tt.todo.Completed)
			assert.True(t, tt.todo.UpdatedAt > originalUpdatedAt)
		})
	}
}

func TestTodo_UpdateTitle(t *testing.T) {
	tests := []struct {
		name     string
		todo     *Todo
		newTitle string
	}{
		{
			name: "updates title",
			todo: &Todo{
				ID:        uuid.New().String(),
				Title:     "Original title",
				Completed: false,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			newTitle: "Updated title",
		},
		{
			name: "updates to empty title",
			todo: &Todo{
				ID:        uuid.New().String(),
				Title:     "Original title",
				Completed: false,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
			newTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := tt.todo.UpdatedAt
			time.Sleep(time.Millisecond) // Ensure time difference

			tt.todo.UpdateTitle(tt.newTitle)

			assert.Equal(t, tt.newTitle, tt.todo.Title)
			assert.True(t, tt.todo.UpdatedAt > originalUpdatedAt)
		})
	}
}
