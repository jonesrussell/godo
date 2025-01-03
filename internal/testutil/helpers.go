// Package testutil provides testing utilities and mock implementations
package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/stretchr/testify/require"
)

// TestingT represents a testing context
type TestingT interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// CreateTestNote creates a note for testing
func CreateTestNote(t TestingT) *note.Note {
	return &note.Note{
		ID:        "test-id",
		Content:   "test content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// RequireNoError asserts that err is nil
func RequireNoError(t TestingT, err error) {
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

// RequireError asserts that err is not nil
func RequireError(t TestingT, err error) {
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// AssertNoteEqual asserts that two notes are equal
func AssertNoteEqual(t *testing.T, expected, actual *note.Note) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Content, actual.Content)
	require.Equal(t, expected.Completed, actual.Completed)
	require.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, time.Second)
	require.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, time.Second)
}

// WithTestContext creates a test context with timeout
func WithTestContext(t TestingT) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
