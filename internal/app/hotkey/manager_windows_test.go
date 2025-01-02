//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

package hotkey

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockQuickNoteService struct {
	mock.Mock
}

func (m *mockQuickNoteService) Show() {
	m.Called()
}

func (m *mockQuickNoteService) Hide() {
	m.Called()
}

func TestPlatformManager_QuickNoteHotkey(t *testing.T) {
	// Create mock quick note service
	quickNote := &mockQuickNoteService{}
	quickNote.On("Show").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager := newPlatformManager(quickNote, binding)

	// Register hotkey
	err := manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Start hotkey listener
	err = manager.Start()
	assert.NoError(t, err, "Should start hotkey listener without error")

	// Give some time for the listener to start
	time.Sleep(100 * time.Millisecond)

	// Clean up
	err = manager.Stop()
	assert.NoError(t, err, "Should stop hotkey listener without error")
}

func TestPlatformManager_InvalidKey(t *testing.T) {
	// Create mock quick note service
	quickNote := &mockQuickNoteService{}

	// Create hotkey binding with invalid key
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "InvalidKey",
	}

	// Create manager
	manager := newPlatformManager(quickNote, binding)

	// Register hotkey should fail
	err := manager.Register()
	assert.Error(t, err, "Should fail to register hotkey with invalid key")
	assert.Contains(t, err.Error(), "unsupported key", "Error should mention unsupported key")
}

func TestPlatformManager_NilQuickNote(t *testing.T) {
	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Should panic when creating manager with nil quick note service
	assert.Panics(t, func() {
		newPlatformManager(nil, binding)
	}, "Should panic when quick note service is nil")
}
