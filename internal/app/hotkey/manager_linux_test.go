//go:build linux && !windows && !darwin
// +build linux,!windows,!darwin

package hotkey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
)

// mockQuickNoteService implements QuickNoteService for testing
type mockQuickNoteService struct {
	shown  bool
	hidden bool
}

func (m *mockQuickNoteService) Show() {
	m.shown = true
	m.hidden = false
}

func (m *mockQuickNoteService) Hide() {
	m.hidden = true
	m.shown = false
}

func TestNewPlatformManager(t *testing.T) {
	mockService := &mockQuickNoteService{}
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Alt"},
		Key:       "N",
	}
	log := logger.NewTestLogger(t)

	manager, err := hotkey.New(mockService, binding, log)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
}

func TestRegisterWithValidBinding(t *testing.T) {
	mockService := &mockQuickNoteService{}
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Alt"},
		Key:       "N",
	}
	log := logger.NewTestLogger(t)

	manager, err := hotkey.New(mockService, binding, log)
	assert.NoError(t, err)

	err = manager.Register()
	assert.NoError(t, err)

	// Clean up
	err = manager.Unregister()
	assert.NoError(t, err)
}

func TestRegisterWithInvalidKey(t *testing.T) {
	mockService := &mockQuickNoteService{}
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Alt"},
		Key:       "InvalidKey",
	}
	log := logger.NewTestLogger(t)

	manager, err := hotkey.New(mockService, binding, log)
	assert.NoError(t, err)

	err = manager.Register()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported key")
}

func TestUnregisterWithoutRegistering(t *testing.T) {
	mockService := &mockQuickNoteService{}
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Alt"},
		Key:       "N",
	}
	log := logger.NewTestLogger(t)

	manager, err := hotkey.New(mockService, binding, log)
	assert.NoError(t, err)

	err = manager.Unregister()
	assert.NoError(t, err)
}

func TestModifierConversion(t *testing.T) {
	testCases := []struct {
		name      string
		modifiers []string
		key       string
		wantErr   bool
	}{
		{
			name:      "All supported modifiers",
			modifiers: []string{"Ctrl", "Alt", "Shift"},
			key:       "N",
			wantErr:   false,
		},
		{
			name:      "Only Ctrl",
			modifiers: []string{"Ctrl"},
			key:       "N",
			wantErr:   false,
		},
		{
			name:      "Only Alt",
			modifiers: []string{"Alt"},
			key:       "N",
			wantErr:   false,
		},
		{
			name:      "Only Shift",
			modifiers: []string{"Shift"},
			key:       "N",
			wantErr:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockQuickNoteService{}
			binding := &common.HotkeyBinding{
				Modifiers: tc.modifiers,
				Key:       tc.key,
			}
			log := logger.NewTestLogger(t)

			manager, err := hotkey.New(mockService, binding, log)
			assert.NoError(t, err)

			err = manager.Register()

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				err = manager.Unregister()
				assert.NoError(t, err)
			}
		})
	}
}
