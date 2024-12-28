//go:build docker

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
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

	if window == nil {
		t.Error("newWindow() returned nil")
	}

	dockerWin, ok := window.(*dockerWindow)
	if !ok {
		t.Error("newWindow() did not return a *dockerWindow")
	}

	if dockerWin.store != store {
		t.Error("store was not properly set")
	}
}

func TestDockerWindowInitialize(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	app := test.NewApp()
	log := &mockLogger{}

	window.Initialize(app, log)

	dockerWin := window.(*dockerWindow)
	if dockerWin.log != log {
		t.Error("logger was not properly set")
	}
}

func TestDockerWindowShowHide(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	app := test.NewApp()
	log := &mockLogger{}

	window.Initialize(app, log)

	// These are no-op functions in Docker, but we should test that they don't panic
	window.Show()
	window.Hide()
}
