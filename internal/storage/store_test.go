package storage

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MemoryStore)
		run     func(*MemoryStore) error
		check   func(*testing.T, *MemoryStore, error)
		cleanup func(*MemoryStore)
	}{
		{
			name:  "add and retrieve task",
			setup: func(s *MemoryStore) {},
			run: func(s *MemoryStore) error {
				task := Task{
					ID:        "test-1",
					Content:   "Test Task",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				return s.Add(task)
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Len(t, tasks, 1)
				assert.Equal(t, "test-1", tasks[0].ID)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name: "update existing task",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Original Content",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.Add(task)
			},
			run: func(s *MemoryStore) error {
				task := Task{
					ID:        "test-1",
					Content:   "Updated Content",
					Done:      true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				return s.Update(task)
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Len(t, tasks, 1)
				assert.Equal(t, "Updated Content", tasks[0].Content)
				assert.True(t, tasks[0].Done)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name:  "update non-existent task",
			setup: func(s *MemoryStore) {},
			run: func(s *MemoryStore) error {
				task := Task{
					ID:        "non-existent",
					Content:   "Content",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				return s.Update(task)
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.ErrorIs(t, err, ErrTaskNotFound)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name: "delete existing task",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Test Task",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.Add(task)
			},
			run: func(s *MemoryStore) error {
				return s.Delete("test-1")
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name:  "delete non-existent task",
			setup: func(s *MemoryStore) {},
			run: func(s *MemoryStore) error {
				return s.Delete("non-existent")
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.ErrorIs(t, err, ErrTaskNotFound)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name: "get task by ID",
			setup: func(s *MemoryStore) {
				task := Task{
					ID:        "test-1",
					Content:   "Test Task",
					Done:      false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				s.Add(task)
			},
			run: func(s *MemoryStore) error {
				task, err := s.GetByID("test-1")
				if err != nil {
					return err
				}
				if task == nil {
					return ErrTaskNotFound
				}
				return nil
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				task, err := s.GetByID("test-1")
				assert.NoError(t, err)
				assert.NotNil(t, task)
				assert.Equal(t, "test-1", task.ID)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name:  "get non-existent task by ID",
			setup: func(s *MemoryStore) {},
			run: func(s *MemoryStore) error {
				task, err := s.GetByID("non-existent")
				if err != nil {
					return err
				}
				if task == nil {
					return ErrTaskNotFound
				}
				return nil
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.ErrorIs(t, err, ErrTaskNotFound)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name:  "list empty store",
			setup: func(s *MemoryStore) {},
			run: func(s *MemoryStore) error {
				_, err := s.List()
				return err
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Empty(t, tasks)
			},
			cleanup: func(s *MemoryStore) {},
		},
		{
			name: "list multiple tasks",
			setup: func(s *MemoryStore) {
				tasks := []Task{
					{
						ID:        "test-1",
						Content:   "Task 1",
						Done:      false,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        "test-2",
						Content:   "Task 2",
						Done:      true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				for _, task := range tasks {
					s.Add(task)
				}
			},
			run: func(s *MemoryStore) error {
				_, err := s.List()
				return err
			},
			check: func(t *testing.T, s *MemoryStore, err error) {
				assert.NoError(t, err)
				tasks, err := s.List()
				assert.NoError(t, err)
				assert.Len(t, tasks, 2)
				assert.Equal(t, "test-1", tasks[0].ID)
				assert.Equal(t, "test-2", tasks[1].ID)
			},
			cleanup: func(s *MemoryStore) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore()
			tt.setup(store)
			err := tt.run(store)
			tt.check(t, store, err)
			tt.cleanup(store)
		})
	}
}

func TestMemoryStoreConcurrent(t *testing.T) {
	store := NewMemoryStore()
	const numGoroutines = 10
	const numOperations = 100

	// Add tasks concurrently
	t.Run("concurrent adds", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					task := Task{
						ID:        fmt.Sprintf("task-%d-%d", routineID, j),
						Content:   fmt.Sprintf("Task %d-%d", routineID, j),
						Done:      false,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					err := store.Add(task)
					assert.NoError(t, err)
				}
			}(i)
		}
		wg.Wait()

		// Verify all tasks were added
		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, numGoroutines*numOperations)
	})

	// Update tasks concurrently
	t.Run("concurrent updates", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					task := Task{
						ID:        fmt.Sprintf("task-%d-%d", routineID, j),
						Content:   fmt.Sprintf("Updated Task %d-%d", routineID, j),
						Done:      true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					err := store.Update(task)
					assert.NoError(t, err)
				}
			}(i)
		}
		wg.Wait()

		// Verify all tasks were updated
		tasks, err := store.List()
		assert.NoError(t, err)
		for _, task := range tasks {
			assert.True(t, task.Done)
			assert.Contains(t, task.Content, "Updated Task")
		}
	})

	// Delete tasks concurrently
	t.Run("concurrent deletes", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					err := store.Delete(fmt.Sprintf("task-%d-%d", routineID, j))
					assert.NoError(t, err)
				}
			}(i)
		}
		wg.Wait()

		// Verify all tasks were deleted
		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})
}
