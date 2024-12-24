//go:build !docker
// +build !docker

package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockUI implements a mock UI for testing
type MockUI struct {
	ShowCalled  bool
	SetupCalled bool
	HideCalled  bool
}

func (m *MockUI) Show() {
	m.ShowCalled = true
}

func (m *MockUI) Setup() {
	m.SetupCalled = true
}

func (m *MockUI) Hide() {
	m.HideCalled = true
}

// MockApplication implements a mock application for testing
type MockApplication struct {
	RunCalled        bool
	SetupUICalled    bool
	GetVersionCalled bool
	ReturnVersion    string
	ReturnError      error
}

func (m *MockApplication) Run() error {
	m.RunCalled = true
	return m.ReturnError
}

func (m *MockApplication) SetupUI() {
	m.SetupUICalled = true
}

func (m *MockApplication) GetVersion() string {
	m.GetVersionCalled = true
	return m.ReturnVersion
}

func TestMockUI(t *testing.T) {
	t.Run("Show sets ShowCalled", func(t *testing.T) {
		ui := &MockUI{}
		ui.Show()
		assert.True(t, ui.ShowCalled)
	})

	t.Run("Setup sets SetupCalled", func(t *testing.T) {
		ui := &MockUI{}
		ui.Setup()
		assert.True(t, ui.SetupCalled)
	})

	t.Run("Hide sets HideCalled", func(t *testing.T) {
		ui := &MockUI{}
		ui.Hide()
		assert.True(t, ui.HideCalled)
	})
}

func TestMockApplication(t *testing.T) {
	t.Run("Run sets RunCalled and returns error", func(t *testing.T) {
		app := &MockApplication{ReturnError: errors.New("test error")}
		err := app.Run()
		assert.True(t, app.RunCalled)
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("SetupUI sets SetupUICalled", func(t *testing.T) {
		app := &MockApplication{}
		app.SetupUI()
		assert.True(t, app.SetupUICalled)
	})

	t.Run("GetVersion sets GetVersionCalled and returns version", func(t *testing.T) {
		app := &MockApplication{ReturnVersion: "1.0.0"}
		version := app.GetVersion()
		assert.True(t, app.GetVersionCalled)
		assert.Equal(t, "1.0.0", version)
	})
}
