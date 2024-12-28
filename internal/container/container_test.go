//go:build !docker && wireinject && windows

package container

import (
	"os"
	"testing"

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
	opts := &LoggerOptions{
		Level:       ProvideLogLevel(),
		Output:      ProvideLogOutputPaths(),
		ErrorOutput: ProvideErrorOutputPaths(),
	}
	logger, cleanup, err := ProvideLogger(opts)
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
	opts := &HotkeyOptions{
		Key:       ProvideKeyCode(),
		Modifiers: ProvideModifierKeys(),
	}
	manager, err := ProvideHotkeyManager(opts)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
}

func TestProvideHTTPConfig(t *testing.T) {
	opts := &HTTPOptions{
		Port:              ProvideHTTPPort(),
		ReadTimeout:       ProvideReadTimeout(),
		WriteTimeout:      ProvideWriteTimeout(),
		ReadHeaderTimeout: ProvideHeaderTimeout(),
		IdleTimeout:       ProvideIdleTimeout(),
	}
	config := ProvideHTTPConfig(opts)
	assert.NotNil(t, config)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, 30, config.ReadTimeout)
	assert.Equal(t, 30, config.WriteTimeout)
	assert.Equal(t, 10, config.ReadHeaderTimeout)
	assert.Equal(t, 120, config.IdleTimeout)
}
