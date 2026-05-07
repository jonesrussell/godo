package config_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/config"
)

func TestNewDefaultConfig_APIInsecureSkipVerifyFalse(t *testing.T) {
	t.Parallel()
	cfg := config.NewDefaultConfig()
	if cfg.Storage.API.InsecureSkipVerify {
		t.Fatalf("expected InsecureSkipVerify false by default, got true")
	}
}

func TestValidateConfig_RejectsCredentialsInAPIBaseURL(t *testing.T) {
	t.Parallel()
	cfg := config.NewDefaultConfig()
	cfg.App.Name = "Godo"
	cfg.Logger.Level = "info"
	cfg.HTTP.StartupTimeout = 5
	cfg.HTTP.ShutdownTimeout = 5
	cfg.Storage.API.BaseURL = "https://user:secret@api.example.com/v1"

	if err := config.ValidateConfig(cfg); err == nil {
		t.Fatal("expected validation error for credentials in base URL")
	}
}
