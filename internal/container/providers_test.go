package container

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{
			name:     "returns development when GODO_ENV is not set",
			envVar:   "",
			expected: "development",
		},
		{
			name:     "returns GODO_ENV value when set",
			envVar:   "production",
			expected: "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.envVar != "" {
				os.Setenv("GODO_ENV", tt.envVar)
				defer os.Unsetenv("GODO_ENV")
			} else {
				os.Unsetenv("GODO_ENV")
			}

			// Test
			result := provideEnvironment()

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}
