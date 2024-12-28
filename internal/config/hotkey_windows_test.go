//go:build windows

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowsGetDefaultQuickNoteModifiers(t *testing.T) {
	modifiers := GetDefaultQuickNoteModifiers()
	assert.NotNil(t, modifiers, "Default quick note modifiers should not be nil")
	assert.NotEmpty(t, modifiers, "Default quick note modifiers should not be empty")
}
