package hotkey

// MockManager implements Manager interface for testing
type MockManager struct {
	RegisterCalled   bool
	UnregisterCalled bool
	StartCalled      bool
	StopCalled       bool
}

// NewMockManager creates a new mock hotkey manager
func NewMockManager() *MockManager {
	return &MockManager{}
}

func (m *MockManager) Register() error {
	m.RegisterCalled = true
	return nil
}

func (m *MockManager) Unregister() error {
	m.UnregisterCalled = true
	return nil
}

func (m *MockManager) Start() error {
	m.StartCalled = true
	return nil
}

func (m *MockManager) Stop() error {
	m.StopCalled = true
	return nil
}
