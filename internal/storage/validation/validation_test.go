package validation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/validation"
)

type mockTaskReader struct {
	storage.TaskReader
}

func TestTaskValidator_ValidateTask(t *testing.T) {
	validator := validation.NewTaskValidator(&mockTaskReader{})

	tests := []struct {
		name    string
		task    storage.Task
		wantErr bool
		errType interface{}
		field   string
	}{
		{
			name: "valid task",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Test content",
				Done:      false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty ID",
			task: storage.Task{
				ID:        "",
				Content:   "Test content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errType: &storage.ValidationError{},
			field:   "id",
		},
		{
			name: "empty content",
			task: storage.Task{
				ID:        "test-1",
				Content:   "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errType: &storage.ValidationError{},
			field:   "content",
		},
		{
			name: "zero created_at",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Test content",
				CreatedAt: time.Time{},
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errType: &storage.ValidationError{},
			field:   "timestamps",
		},
		{
			name: "zero updated_at",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Test content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Time{},
			},
			wantErr: true,
			errType: &storage.ValidationError{},
			field:   "timestamps",
		},
		{
			name: "updated_at before created_at",
			task: storage.Task{
				ID:        "test-1",
				Content:   "Test content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
			wantErr: true,
			errType: &storage.ValidationError{},
			field:   "timestamps",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTask(tt.task)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.IsType(t, tt.errType, err)
				}
				if tt.field != "" {
					validationErr, ok := err.(*storage.ValidationError)
					require.True(t, ok)
					assert.Equal(t, tt.field, validationErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateConnection(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "no error",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "connection error",
			err:     errors.New("connection failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateConnection(tt.err)
			if tt.wantErr {
				require.Error(t, err)
				assert.IsType(t, &storage.ConnectionError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTransaction(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr bool
	}{
		{
			name:    "no error",
			err:     nil,
			wantErr: false,
		},
		{
			name:    "transaction error",
			err:     errors.New("transaction failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateTransaction(tt.err)
			if tt.wantErr {
				require.Error(t, err)
				assert.IsType(t, &storage.TransactionError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
