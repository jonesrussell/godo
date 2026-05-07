package api

import (
	"crypto/tls"

	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
)

// tlsClientConfigForStorage builds the TLS client configuration for API-backed storage.
// InsecureSkipVerify remains false unless TLSInsecureSkipVerify is explicitly set on the config.
func tlsClientConfigForStorage(cfg domainstorage.APIConfig) *tls.Config {
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.TLSInsecureSkipVerify,
	}
}
