//go:build !docker

package container

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	manager := ProvideHotkeyManager()
	assert.NotNil(t, manager)
}
