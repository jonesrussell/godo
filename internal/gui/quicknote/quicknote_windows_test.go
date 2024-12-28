//go:build windows && !linux

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowsWindowImplementation(t *testing.T) {
	store := &mockStore{}
	log := &mockLogger{}
	window := New(store, log)
	require.NotNil(t, window)

	// Test initialization
	app := test.NewApp()
	defer app.Quit()

	// Test window operations
	window.Show()
	window.Hide()
}

func TestWindowsWindowInterface(t *testing.T) {
	store := &mockStore{}
	log := &mockLogger{}
	window := New(store, log)
	require.NotNil(t, window)

	// Test that the window implements the Interface correctly
	app := test.NewApp()
	defer app.Quit()
	window.Show()
	window.Hide()
}

func TestWindowsWindowWithNilStore(t *testing.T) {
	// Test that the window can be created with a nil store
	log := &mockLogger{}
	window := New(nil, log)
	require.NotNil(t, window)

	// Operations should not panic with nil store
	app := test.NewApp()
	defer app.Quit()
	window.Show()
	window.Hide()
}

func TestWindowsWindowWithNilLogger(t *testing.T) {
	store := &mockStore{}
	window := New(store, nil)
	require.NotNil(t, window)

	// Test that initialization with nil logger doesn't panic
	app := test.NewApp()
	defer app.Quit()
	window.Show()
	window.Hide()
}

func TestWindowsWindowTaskOperations(t *testing.T) {
	store := &mockStore{}
	log := &mockLogger{}
	window := New(store, log)
	require.NotNil(t, window)

	app := test.NewApp()
	defer app.Quit()

	// Test task creation
	store.tasks = []storage.Task{
		{ID: "1", Content: "Task 1"},
		{ID: "2", Content: "Task 2"},
	}
	window.Show()
	assert.True(t, store.listCalled)

	// Test task update
	task := storage.Task{ID: "1", Content: "Updated Task 1", Done: true}
	store.tasks = []storage.Task{task}
	window.Show()
	assert.True(t, store.listCalled)
}

func TestWindowsWindowErrorHandling(t *testing.T) {
	store := &mockStore{err: assert.AnError}
	log := &mockLogger{}
	window := New(store, log)
	require.NotNil(t, window)

	app := test.NewApp()
	defer app.Quit()

	// Test error handling during task operations
	window.Show()
	assert.True(t, store.listCalled)
	assert.True(t, log.errorCalled)
}

type mockLogger struct {
	debugCalled bool
	infoCalled  bool
	warnCalled  bool
	errorCalled bool
}

func (m *mockLogger) Debug(msg string, keysAndValues ...interface{})         { m.debugCalled = true }
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})          { m.infoCalled = true }
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})          { m.warnCalled = true }
func (m *mockLogger) Error(msg string, keysAndValues ...interface{})         { m.errorCalled = true }
func (m *mockLogger) WithError(err error) logger.Logger                      { return m }
func (m *mockLogger) WithField(key string, value interface{}) logger.Logger  { return m }
func (m *mockLogger) WithFields(fields map[string]interface{}) logger.Logger { return m }
