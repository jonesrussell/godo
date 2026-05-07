package api

import (
	"testing"

	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
)

func TestTLSClientConfigForStorage_DefaultSecure(t *testing.T) {
	t.Parallel()
	cfg := domainstorage.APIConfig{
		BaseURL:               "https://example.com/api",
		TLSInsecureSkipVerify: false,
	}
	tc := tlsClientConfigForStorage(cfg)
	if tc.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify false by default")
	}
}

func TestTLSClientConfigForStorage_ExplicitInsecureOptIn(t *testing.T) {
	t.Parallel()
	cfg := domainstorage.APIConfig{
		BaseURL:               "https://example.com/api",
		TLSInsecureSkipVerify: true,
	}
	tc := tlsClientConfigForStorage(cfg)
	if !tc.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify true when opted in")
	}
}
