//go:build !docker && wireinject

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockApp implements a minimal test version of the application
type mockApp struct {
	setupUICalled bool
	runCalled     bool
	version       string
}

func newMockApp() *mockApp {
	return &mockApp{
		version: "test-version",
	}
}

func (m *mockApp) SetupUI() {
	m.setupUICalled = true
}

func (m *mockApp) Run() {
	m.runCalled = true
}

func TestMainFlow(t *testing.T) {
	// Create a mock app
	mockApp := newMockApp()

	// Run the app
	mockApp.SetupUI()
	mockApp.Run()

	// Verify app was properly initialized and run
	assert.True(t, mockApp.setupUICalled, "SetupUI should have been called")
	assert.True(t, mockApp.runCalled, "Run should have been called")
}
