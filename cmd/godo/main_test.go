//go:build !docker && wireinject

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockApp implements a minimal test version of the application
type mockApp struct {
	setupUICalled   bool
	runCalled       bool
	version         string
	shouldFailSetup bool
	shouldFailRun   bool
	closed          bool
}

func newMockApp() *mockApp {
	return &mockApp{
		version: "test-version",
	}
}

func (m *mockApp) SetupUI() error {
	if m.closed {
		return errors.New("app is closed")
	}
	if m.shouldFailSetup {
		return errors.New("setup failed")
	}
	m.setupUICalled = true
	return nil
}

func (m *mockApp) Run() error {
	if m.closed {
		return errors.New("app is closed")
	}
	if m.shouldFailRun {
		return errors.New("run failed")
	}
	m.runCalled = true
	return nil
}

func (m *mockApp) Close() error {
	if m.closed {
		return errors.New("app already closed")
	}
	m.closed = true
	return nil
}

func TestMainFlow(t *testing.T) {
	t.Run("successful flow", func(t *testing.T) {
		mockApp := newMockApp()

		err := mockApp.SetupUI()
		require.NoError(t, err)

		err = mockApp.Run()
		require.NoError(t, err)

		assert.True(t, mockApp.setupUICalled, "SetupUI should have been called")
		assert.True(t, mockApp.runCalled, "Run should have been called")
	})

	t.Run("setup failure", func(t *testing.T) {
		mockApp := newMockApp()
		mockApp.shouldFailSetup = true

		err := mockApp.SetupUI()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "setup failed")
		assert.False(t, mockApp.setupUICalled)
	})

	t.Run("run failure", func(t *testing.T) {
		mockApp := newMockApp()
		mockApp.shouldFailRun = true

		err := mockApp.SetupUI()
		require.NoError(t, err)

		err = mockApp.Run()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "run failed")
		assert.True(t, mockApp.setupUICalled)
		assert.False(t, mockApp.runCalled)
	})

	t.Run("closed app", func(t *testing.T) {
		mockApp := newMockApp()

		err := mockApp.Close()
		require.NoError(t, err)

		err = mockApp.SetupUI()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app is closed")

		err = mockApp.Run()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app is closed")
	})

	t.Run("double close", func(t *testing.T) {
		mockApp := newMockApp()

		err := mockApp.Close()
		require.NoError(t, err)

		err = mockApp.Close()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "app already closed")
	})
}
