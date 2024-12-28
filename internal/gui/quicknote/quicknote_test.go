package quicknote

import (
	"testing"

	"github.com/jonesrussell/godo/internal/gui"
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

type mockQuickNote struct {
	gui.QuickNote
	showCalled bool
	hideCalled bool
}

func (m *mockQuickNote) Show() {
	m.showCalled = true
}

func (m *mockQuickNote) Hide() {
	m.hideCalled = true
}

func TestQuickNoteShowHide(t *testing.T) {
	qn := &mockQuickNote{}

	// Test Show
	qn.Show()
	assert.True(t, qn.showCalled)

	// Test Hide
	qn.Hide()
	assert.True(t, qn.hideCalled)
}

func TestQuickNoteInterface(t *testing.T) {
	var quickNote gui.QuickNote = &mockQuickNote{}
	require.NotNil(t, quickNote)

	quickNote.Show()
	quickNote.Hide()
}
