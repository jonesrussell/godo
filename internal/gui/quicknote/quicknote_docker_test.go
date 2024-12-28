//go:build docker && !windows
// +build docker,!windows

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStore struct {
	storage.Store
	addCalled    bool
	updateCalled bool
	deleteCalled bool
	listCalled   bool
	tasks        []storage.Task
	err          error
}

type mockLogger struct {
	debugCalled bool
	infoCalled  bool
	warnCalled  bool
	errorCalled bool
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

func TestDockerWindowImplementation(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	require.NotNil(t, window)

	// Test initialization
	app := test.NewApp()
	defer app.Quit()
	log := &mockLogger{}
	window.Initialize(app, log)

	// Test that Show and Hide are no-ops (don't panic)
	window.Show()
	window.Hide()
}

func TestDockerWindowInterface(t *testing.T) {
	store := &mockStore{}
	var window Interface = newWindow(store)
	require.NotNil(t, window)

	// Test that the window implements the Interface correctly
	app := test.NewApp()
	defer app.Quit()
	log := &mockLogger{}
	window.Initialize(app, log)
	window.Show()
	window.Hide()
}

func TestDockerWindowWithNilStore(t *testing.T) {
	// Test that the window can be created with a nil store
	window := newWindow(nil)
	require.NotNil(t, window)

	// Operations should not panic with nil store
	app := test.NewApp()
	defer app.Quit()
	log := &mockLogger{}
	window.Initialize(app, log)
	window.Show()
	window.Hide()
}

func TestDockerWindowWithNilLogger(t *testing.T) {
	store := &mockStore{}
	window := newWindow(store)
	require.NotNil(t, window)

	// Test that initialization with nil logger doesn't panic
	app := test.NewApp()
	defer app.Quit()
	window.Initialize(app, nil)
	window.Show()
	window.Hide()
}

func (m *mockLogger) Debug(msg string, keysAndValues ...interface{}) { m.debugCalled = true }
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})  { m.infoCalled = true }
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})  { m.warnCalled = true }
func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) { m.errorCalled = true }

func (m *mockStore) Add(task storage.Task) error {
	m.addCalled = true
	if m.err != nil {
		return m.err
	}
	m.tasks = append(m.tasks, task)
	return nil
}

func (m *mockStore) Update(task storage.Task) error {
	m.updateCalled = true
	if m.err != nil {
		return m.err
	}
	for i, t := range m.tasks {
		if t.ID == task.ID {
			m.tasks[i] = task
			return nil
		}
	}
	return storage.ErrTaskNotFound
}

func (m *mockStore) Delete(id string) error {
	m.deleteCalled = true
	if m.err != nil {
		return m.err
	}
	for i, task := range m.tasks {
		if task.ID == id {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return storage.ErrTaskNotFound
}

func (m *mockStore) List() ([]storage.Task, error) {
	m.listCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return m.tasks, nil
}

func (m *mockStore) GetByID(id string) (*storage.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, task := range m.tasks {
		if task.ID == id {
			return &task, nil
		}
	}
	return nil, storage.ErrTaskNotFound
}

func (m *mockStore) Close() error {
	return nil
}
