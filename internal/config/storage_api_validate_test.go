package config

import "testing"

func TestValidateConfig_APIBaseURL(t *testing.T) {
	t.Parallel()
	cfg := NewDefaultConfig()
	cfg.Storage.Type = "api"
	cfg.Storage.API.BaseURL = "https://api.example.com/v1"
	if err := ValidateConfig(cfg); err != nil {
		t.Fatalf("expected ok: %v", err)
	}
}

func TestValidateConfig_APIBaseURL_RejectsUserinfo(t *testing.T) {
	t.Parallel()
	cfg := NewDefaultConfig()
	cfg.Storage.Type = "api"
	cfg.Storage.API.BaseURL = "https://user:pass@example.com/api"
	if err := ValidateConfig(cfg); err == nil {
		t.Fatal("expected validation error")
	}
}
