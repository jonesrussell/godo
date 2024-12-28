//go:build !docker

package container

import (
	"os"
	"testing"

	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/stretchr/testify/assert"
)

// mockQuickNote is a mock implementation of quicknote.Interface for testing
type mockQuickNote struct{}

// Ensure mockQuickNote implements quicknote.Interface
var _ quicknote.Interface = (*mockQuickNote)(nil)

func (m *mockQuickNote) Setup() error { return nil }
func (m *mockQuickNote) Show()        {}
func (m *mockQuickNote) Hide()        {}

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "true" {
		os.Exit(0) // Skip tests in CI environment
	}
	os.Exit(m.Run())
}

func TestProvideLogger(t *testing.T) {
	logger, cleanup, err := ProvideLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, cleanup)
	cleanup()
}

func TestProvideFyneApp(t *testing.T) {
	app := ProvideFyneApp()
	assert.NotNil(t, app)
}

func TestProvideStorage(t *testing.T) {
	store := ProvideStorage()
	assert.NotNil(t, store)
}

func TestProvideHotkeyManager(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping hotkey test in CI environment")
	}

	// Create a mock quick note service
	mockNote := &mockQuickNote{}

	manager := ProvideHotkeyManager(mockNote)
	assert.NotNil(t, manager)
}
