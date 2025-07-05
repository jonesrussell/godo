package common_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/jonesrussell/godo/internal/common"
)

func TestLoadConfig(t *testing.T) {
	yamlConfig := `
hotkeys:
  quick_note:
    modifiers: ["Ctrl", "Alt"]
    key: "G"
http:
  port: 8080
  read_timeout: 30
  write_timeout: 30
  read_header_timeout: 10
  idle_timeout: 120
logger:
  level: "info"
  console: true
  file: false
  file_path: ""
  output: ["stdout"]
  error_output: ["stderr"]
`
	var config common.Config
	err := yaml.Unmarshal([]byte(yamlConfig), &config)
	require.NoError(t, err)

	// Test hotkey configuration
	assert.Equal(t, []string{"Ctrl", "Alt"}, config.Hotkeys.QuickNote.Modifiers)
	assert.Equal(t, "G", config.Hotkeys.QuickNote.Key)

	// Test HTTP configuration
	assert.Equal(t, 8080, config.HTTP.Port)
	assert.Equal(t, 30, config.HTTP.ReadTimeout)
	assert.Equal(t, 30, config.HTTP.WriteTimeout)
	assert.Equal(t, 10, config.HTTP.ReadHeaderTimeout)
	assert.Equal(t, 120, config.HTTP.IdleTimeout)

	// Test logger configuration
	assert.Equal(t, "info", config.Logger.Level)
	assert.True(t, config.Logger.Console)
	assert.False(t, config.Logger.File)
	assert.Empty(t, config.Logger.FilePath)
	assert.Equal(t, []string{"stdout"}, config.Logger.Output)
	assert.Equal(t, []string{"stderr"}, config.Logger.ErrorOutput)
}

func TestHTTPConfigTimeouts(t *testing.T) {
	config := common.HTTPConfig{
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       120,
	}

	// Test timeout conversion methods
	assert.Equal(t, 30*time.Second, config.GetReadTimeout())
	assert.Equal(t, 30*time.Second, config.GetWriteTimeout())
	assert.Equal(t, 10*time.Second, config.GetReadHeaderTimeout())
	assert.Equal(t, 120*time.Second, config.GetIdleTimeout())
}

func TestHotkeyBindingValidation(t *testing.T) {
	testCases := []struct {
		name      string
		binding   common.HotkeyBinding
		wantError bool
	}{
		{
			name: "Valid binding",
			binding: common.HotkeyBinding{
				Modifiers: []string{"Ctrl", "Alt"},
				Key:       "G",
			},
			wantError: false,
		},
		{
			name: "Empty modifiers",
			binding: common.HotkeyBinding{
				Modifiers: []string{},
				Key:       "G",
			},
			wantError: true,
		},
		{
			name: "Empty key",
			binding: common.HotkeyBinding{
				Modifiers: []string{"Ctrl", "Alt"},
				Key:       "",
			},
			wantError: true,
		},
		{
			name: "Invalid modifier",
			binding: common.HotkeyBinding{
				Modifiers: []string{"Invalid", "Alt"},
				Key:       "G",
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.binding.Validate()
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogConfigValidation(t *testing.T) {
	testCases := []struct {
		name      string
		config    common.LogConfig
		wantError bool
	}{
		{
			name: "Valid config",
			config: common.LogConfig{
				Level:       "info",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantError: false,
		},
		{
			name: "Invalid level",
			config: common.LogConfig{
				Level:       "invalid",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantError: true,
		},
		{
			name: "Empty output",
			config: common.LogConfig{
				Level:       "info",
				Output:      []string{},
				ErrorOutput: []string{"stderr"},
			},
			wantError: true,
		},
		{
			name: "Empty error output",
			config: common.LogConfig{
				Level:       "info",
				Output:      []string{"stdout"},
				ErrorOutput: []string{},
			},
			wantError: true,
		},
		{
			name: "File enabled without path",
			config: common.LogConfig{
				Level:       "info",
				File:        true,
				FilePath:    "",
				Output:      []string{"stdout"},
				ErrorOutput: []string{"stderr"},
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHTTPConfigValidation(t *testing.T) {
	testCases := []struct {
		name      string
		config    common.HTTPConfig
		wantError bool
	}{
		{
			name: "Valid config",
			config: common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       30,
				WriteTimeout:      30,
				ReadHeaderTimeout: 10,
				IdleTimeout:       120,
			},
			wantError: false,
		},
		{
			name: "Invalid port",
			config: common.HTTPConfig{
				Port:              0,
				ReadTimeout:       30,
				WriteTimeout:      30,
				ReadHeaderTimeout: 10,
				IdleTimeout:       120,
			},
			wantError: true,
		},
		{
			name: "Invalid read timeout",
			config: common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       0,
				WriteTimeout:      30,
				ReadHeaderTimeout: 10,
				IdleTimeout:       120,
			},
			wantError: true,
		},
		{
			name: "Invalid write timeout",
			config: common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       30,
				WriteTimeout:      0,
				ReadHeaderTimeout: 10,
				IdleTimeout:       120,
			},
			wantError: true,
		},
		{
			name: "Invalid read header timeout",
			config: common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       30,
				WriteTimeout:      30,
				ReadHeaderTimeout: 0,
				IdleTimeout:       120,
			},
			wantError: true,
		},
		{
			name: "Invalid idle timeout",
			config: common.HTTPConfig{
				Port:              8080,
				ReadTimeout:       30,
				WriteTimeout:      30,
				ReadHeaderTimeout: 10,
				IdleTimeout:       0,
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
