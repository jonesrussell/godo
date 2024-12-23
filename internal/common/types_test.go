package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHotkeyBinding_String(t *testing.T) {
	hk := HotkeyBinding{
		Modifiers: []string{"ctrl", "shift"},
		Key:       "n",
	}
	assert.Equal(t, "n", hk.String())
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "with wrapped error",
			err: &Error{
				Code:    "TEST_ERROR",
				Message: "test message",
				Err:     errors.New("wrapped error"),
			},
			expected: "TEST_ERROR: test message: wrapped error",
		},
		{
			name: "without wrapped error",
			err: &Error{
				Code:    "TEST_ERROR",
				Message: "test message",
			},
			expected: "TEST_ERROR: test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestNewError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := NewError("TEST_ERROR", "test message", wrappedErr)

	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, wrappedErr, err.Err)
}
