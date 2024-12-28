//go:build !docker

package container

import (
	"os"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/stretchr/testify/assert"
)

// mockQuickNote is a mock implementation of gui.QuickNote for testing
type mockQuickNote struct{}

// Ensure mockQuickNote implements gui.QuickNote
var _ gui.QuickNote = (*mockQuickNote)(nil)

func (m *mockQuickNote) Show() {}
func (m *mockQuickNote) Hide() {}

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

func TestProvideHotkeyManager(t *testing.T) {
	binding := &common.HotkeyBinding{
		Key:       "N",
		Modifiers: []string{"Ctrl"},
	}
	manager, err := ProvideHotkeyManager(binding)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
}
