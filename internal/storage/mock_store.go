package storage

// MockStore implements Store interface for testing
type MockStore struct {
	SaveCalled   bool
	LoadCalled   bool
	CloseCalled  bool
	DeleteCalled bool
	AddCalled    bool
	Data         map[string]interface{}
}

// NewMockStore creates a new mock store
func NewMockStore() *MockStore {
	return &MockStore{
		Data: make(map[string]interface{}),
	}
}

func (m *MockStore) Save(key string, value interface{}) error {
	m.SaveCalled = true
	m.Data[key] = value
	return nil
}

func (m *MockStore) Load(key string) (interface{}, error) {
	m.LoadCalled = true
	return m.Data[key], nil
}

func (m *MockStore) Delete(key string) error {
	m.DeleteCalled = true
	delete(m.Data, key)
	return nil
}

func (m *MockStore) Close() error {
	m.CloseCalled = true
	return nil
}

func (m *MockStore) Add(key string, value interface{}) error {
	m.AddCalled = true
	m.Data[key] = value
	return nil
}
