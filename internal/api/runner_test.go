package api

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunnerLifecycle(t *testing.T) {
	// Setup
	log := logger.NewTestLogger(t)
	store := testutil.NewMockStore()
	config := &common.HTTPConfig{
		Port:              0, // Let OS choose port
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       120,
	}

	// Create runner
	runner := NewRunner(store, log, config)
	require.NotNil(t, runner)

	// Test start
	runner.Start(0)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := runner.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestRunnerWithInvalidPort(t *testing.T) {
	// Setup
	log := logger.NewTestLogger(t)
	store := testutil.NewMockStore()
	config := &common.HTTPConfig{
		Port:              -1, // Invalid port
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       120,
	}

	// Create runner
	runner := NewRunner(store, log, config)
	require.NotNil(t, runner)

	// Test start with invalid port
	runner.Start(-1)

	// Give the server a moment to attempt to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := runner.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestRunnerWithClosedStore(t *testing.T) {
	// Setup
	log := logger.NewTestLogger(t)
	store := testutil.NewMockStore()
	config := &common.HTTPConfig{
		Port:              0,
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       120,
	}

	// Close store before starting server
	require.NoError(t, store.Close())

	// Create runner
	runner := NewRunner(store, log, config)
	require.NotNil(t, runner)

	// Test start
	runner.Start(0)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := runner.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestRunnerWithShortTimeout(t *testing.T) {
	// Setup
	log := logger.NewTestLogger(t)
	store := testutil.NewMockStore()
	config := &common.HTTPConfig{
		Port:              0,
		ReadTimeout:       1,
		WriteTimeout:      1,
		ReadHeaderTimeout: 1,
		IdleTimeout:       1,
	}

	// Create runner
	runner := NewRunner(store, log, config)
	require.NotNil(t, runner)

	// Test start
	runner.Start(0)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown with very short context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err := runner.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestRunnerWithNilDependencies(t *testing.T) {
	// Test with nil logger
	t.Run("nil logger", func(t *testing.T) {
		store := testutil.NewMockStore()
		config := &common.HTTPConfig{Port: 0}
		runner := NewRunner(store, nil, config)
		require.NotNil(t, runner)
	})

	// Test with nil config
	t.Run("nil config", func(t *testing.T) {
		store := testutil.NewMockStore()
		log := logger.NewTestLogger(t)
		runner := NewRunner(store, log, nil)
		require.NotNil(t, runner)
	})

	// Test with nil store
	t.Run("nil store", func(t *testing.T) {
		log := logger.NewTestLogger(t)
		config := &common.HTTPConfig{Port: 0}
		runner := NewRunner(nil, log, config)
		require.NotNil(t, runner)
	})
}
