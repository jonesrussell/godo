package mainwindow

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/gui"
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

type mockLogger struct {
	logger.Logger
	debugCalled bool
	infoCalled  bool
	warnCalled  bool
	errorCalled bool
}

func (m *mockLogger) Debug(msg string, keysAndValues ...interface{}) { m.debugCalled = true }
func (m *mockLogger) Info(msg string, keysAndValues ...interface{})  { m.infoCalled = true }
func (m *mockLogger) Warn(msg string, keysAndValues ...interface{})  { m.warnCalled = true }
func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) { m.errorCalled = true }

func TestMainWindowInterface(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{}
	log := &mockLogger{}

	var mainWindow gui.MainWindow = New(store, log)
	require.NotNil(t, mainWindow)

	// Test window operations
	mainWindow.Show()
	mainWindow.Hide()
	mainWindow.CenterOnScreen()

	// Test content setting
	content := test.NewCanvas().Content()
	mainWindow.SetContent(content)

	// Test window resizing
	mainWindow.Resize(fyne.NewSize(800, 600))

	// Test window getter
	window := mainWindow.GetWindow()
	require.NotNil(t, window)
}

func TestMainWindowTaskOperations(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{}
	log := &mockLogger{}
	mainWindow := New(store, log)
	require.NotNil(t, mainWindow)

	// Test task list refresh
	store.tasks = []storage.Task{
		{ID: "1", Content: "Task 1"},
		{ID: "2", Content: "Task 2"},
	}
	mainWindow.Show()
	assert.True(t, store.listCalled)

	// Test task completion
	task := storage.Task{ID: "1", Content: "Task 1", Done: false}
	store.tasks = []storage.Task{task}
	mainWindow.Show()
	assert.True(t, store.listCalled)

	// Test task deletion
	store.tasks = []storage.Task{task}
	mainWindow.Show()
	assert.True(t, store.listCalled)
}

func TestMainWindowErrorHandling(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{err: assert.AnError}
	log := &mockLogger{}
	mainWindow := New(store, log)
	require.NotNil(t, mainWindow)

	// Test error handling during task list refresh
	mainWindow.Show()
	assert.True(t, store.listCalled)
	assert.True(t, log.errorCalled)
}

func TestMainWindowLayout(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{}
	log := &mockLogger{}
	mainWindow := New(store, log)
	require.NotNil(t, mainWindow)

	// Test initial window size
	window := mainWindow.GetWindow()
	require.NotNil(t, window)
	assert.Greater(t, window.Canvas().Size().Width, float32(0))
	assert.Greater(t, window.Canvas().Size().Height, float32(0))

	// Test window resize
	newSize := fyne.NewSize(1000, 800)
	mainWindow.Resize(newSize)
	assert.Equal(t, newSize.Width, window.Canvas().Size().Width)
	assert.Equal(t, newSize.Height, window.Canvas().Size().Height)
}

func TestMainWindowLifecycle(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{}
	log := &mockLogger{}
	mainWindow := New(store, log)
	require.NotNil(t, mainWindow)

	// Test window show/hide cycle
	mainWindow.Show()
	// Note: We can't directly test window visibility in Fyne test mode
	// Instead, we'll verify that the window exists and operations don't panic
	window := mainWindow.GetWindow()
	require.NotNil(t, window)

	mainWindow.Hide()
	// Verify window still exists after hide
	window = mainWindow.GetWindow()
	require.NotNil(t, window)

	// Test window centering
	mainWindow.CenterOnScreen()
	// Note: We can't test the actual position as it depends on screen size
}

func TestMainWindowContentUpdate(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	store := &mockStore{}
	log := &mockLogger{}
	mainWindow := New(store, log)
	require.NotNil(t, mainWindow)

	// Test content update with tasks
	store.tasks = []storage.Task{
		{ID: "1", Content: "Task 1", Done: false},
		{ID: "2", Content: "Task 2", Done: true},
	}
	mainWindow.Show()
	assert.True(t, store.listCalled)

	// Test content update with empty task list
	store.tasks = []storage.Task{}
	mainWindow.Show()
	assert.True(t, store.listCalled)
}
