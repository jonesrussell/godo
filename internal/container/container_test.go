//go:build !docker && wireinject && windows

package container

import (
	"os"
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
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

func TestNew(t *testing.T) {
	// Test successful container creation
	container, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, container)
	assert.NotNil(t, container.App)
	assert.NotNil(t, container.Logger)
	assert.NotNil(t, container.Store)

	// Test app type assertion
	_, ok := container.App.(*app.App)
	assert.True(t, ok)
}

func TestProvideLogger(t *testing.T) {
	tests := []struct {
		name    string
		opts    *LoggerOptions
		wantErr bool
	}{
		{
			name: "valid logger options",
			opts: &LoggerOptions{
				Level:       common.LogLevel("debug"),
				Output:      common.LogOutputPaths{"stdout"},
				ErrorOutput: common.ErrorOutputPaths{"stderr"},
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			opts: &LoggerOptions{
				Level:       common.LogLevel("invalid"),
				Output:      common.LogOutputPaths{"stdout"},
				ErrorOutput: common.ErrorOutputPaths{"stderr"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, cleanup, err := ProvideLogger(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				assert.Nil(t, cleanup)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)
			assert.NotNil(t, cleanup)
			cleanup()
		})
	}
}

func TestProvideHTTPConfig(t *testing.T) {
	tests := []struct {
		name string
		opts *HTTPOptions
		want *common.HTTPConfig
	}{
		{
			name: "default config",
			opts: &HTTPOptions{
				Port:              common.HTTPPort(8080),
				ReadTimeout:       common.ReadTimeoutSeconds(30),
				WriteTimeout:      common.WriteTimeoutSeconds(30),
				ReadHeaderTimeout: common.HeaderTimeoutSeconds(10),
				IdleTimeout:       common.IdleTimeoutSeconds(120),
			},
			want: &common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       30,
				WriteTimeout:      30,
				ReadHeaderTimeout: 10,
				IdleTimeout:       120,
			},
		},
		{
			name: "custom config",
			opts: &HTTPOptions{
				Port:              common.HTTPPort(9090),
				ReadTimeout:       common.ReadTimeoutSeconds(60),
				WriteTimeout:      common.WriteTimeoutSeconds(60),
				ReadHeaderTimeout: common.HeaderTimeoutSeconds(20),
				IdleTimeout:       common.IdleTimeoutSeconds(180),
			},
			want: &common.HTTPConfig{
				Port:              9090,
				ReadTimeout:       60,
				WriteTimeout:      60,
				ReadHeaderTimeout: 20,
				IdleTimeout:       180,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProvideHTTPConfig(tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProvideHotkeyManager(t *testing.T) {
	tests := []struct {
		name    string
		opts    *HotkeyOptions
		wantErr bool
	}{
		{
			name: "valid hotkey options",
			opts: &HotkeyOptions{
				Key:       common.KeyCode("N"),
				Modifiers: common.ModifierKeys{"ctrl", "alt"},
			},
			wantErr: false,
		},
		{
			name: "empty key",
			opts: &HotkeyOptions{
				Key:       common.KeyCode(""),
				Modifiers: common.ModifierKeys{"ctrl", "alt"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := ProvideHotkeyManager(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, manager)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, manager)
		})
	}
}

func TestProvideAppMetadata(t *testing.T) {
	t.Run("app name", func(t *testing.T) {
		name := ProvideAppName()
		assert.Equal(t, "Godo", name.String())
	})

	t.Run("app version", func(t *testing.T) {
		version := ProvideAppVersion()
		assert.NotEmpty(t, version.String())
		assert.Regexp(t, `^\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$`, version.String())
	})

	t.Run("app id", func(t *testing.T) {
		id := ProvideAppID()
		assert.Equal(t, "com.jonesrussell.godo", id.String())
	})
}

func TestProvideDatabaseConfig(t *testing.T) {
	path := ProvideDatabasePath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "godo.db")
}

func TestProvideSQLiteStore(t *testing.T) {
	log := logger.NewTestLogger(t)
	store, err := ProvideSQLiteStore(log)
	assert.NoError(t, err)
	assert.NotNil(t, store)
	defer store.Close()
}

func TestProvideTimeouts(t *testing.T) {
	assert.Equal(t, common.ReadTimeoutSeconds(30), ProvideReadTimeout())
	assert.Equal(t, common.WriteTimeoutSeconds(30), ProvideWriteTimeout())
	assert.Equal(t, common.HeaderTimeoutSeconds(10), ProvideHeaderTimeout())
	assert.Equal(t, common.IdleTimeoutSeconds(120), ProvideIdleTimeout())
}
