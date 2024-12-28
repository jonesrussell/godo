//go:build docker && !windows

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerGetDefaultQuickNoteModifiers(t *testing.T) {
	modifiers := GetDefaultQuickNoteModifiers()
	assert.NotNil(t, modifiers, "Default quick note modifiers should not be nil")
	assert.Empty(t, modifiers, "Default quick note modifiers should be empty in Docker environment")
}
