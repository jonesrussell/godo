package container

import (
	"testing"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviders(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *config.Config
		test    func(*testing.T, *config.Config)
		cleanup func()
	}{
		{
			name: "BootstrapLogger creates basic logger",
			setup: func() *config.Config {
				return nil // BootstrapLogger doesn't need config
			},
			test: func(t *testing.T, _ *config.Config) {
				log, err := BootstrapLogger()
				require.NoError(t, err)
				assert.NotNil(t, log)
			},
		},
		{
			name: "ProvideSQLite creates store with correct path",
			setup: func() *config.Config {
				return &config.Config{
					Database: config.DatabaseConfig{
						Path:         t.TempDir(),
						MaxOpenConns: 1,
						MaxIdleConns: 1,
					},
				}
			},
			test: func(t *testing.T, cfg *config.Config) {
				log, _ := logger.NewZapLogger(&logger.Config{
					Level:   "debug",
					Console: true,
				})
				store, err := ProvideSQLite(cfg, log)
				require.NoError(t, err)
				assert.NotNil(t, store)
				store.Close()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setup()
			tt.test(t, cfg)
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
