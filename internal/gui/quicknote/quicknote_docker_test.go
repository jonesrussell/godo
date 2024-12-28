//go:build docker && !windows
// +build docker,!windows

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	storage.Store
}

type mockLogger struct {
	logger.Logger
}

func TestNewWindow(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)

	assert.NotNil(t, window, "newWindow() should not return nil")

	dockerWin, ok := window.(*dockerWindow)
	assert.True(t, ok, "newWindow() should return a *dockerWindow")
	assert.Equal(t, store, dockerWin.store, "store should be properly set")
}

func TestDockerWindowInitialize(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	app := test.NewApp()
	log := &mockLogger{}

	window.Initialize(app, log)

	dockerWin := window.(*dockerWindow)
	assert.Equal(t, log, dockerWin.log, "logger should be properly set")
}

func TestDockerWindowShowHide(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	app := test.NewApp()
	log := &mockLogger{}

	window.Initialize(app, log)

	// These should be no-op functions in Docker environment
	assert.NotPanics(t, func() {
		window.Show()
		window.Hide()
	}, "Show and Hide should not panic in Docker environment")
}
