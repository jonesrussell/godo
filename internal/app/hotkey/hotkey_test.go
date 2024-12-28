//go:build !ci && !android && !ios && !wasm && !test_web_driver && !docker

package hotkey

// TestManager is a mock hotkey implementation for testing
type TestManager struct {
	quickNote QuickNoteService
	isActive  bool
}

// NewTestManager creates a new test hotkey manager
func NewTestManager(quickNote QuickNoteService) Manager {
	return &TestManager{
		quickNote: quickNote,
	}
}

func (h *TestManager) Register() error {
	h.isActive = true
	return nil
}

func (h *TestManager) Unregister() error {
	h.isActive = false
	return nil
}

// Trigger simulates a hotkey press
func (h *TestManager) Trigger() {
	if h.isActive && h.quickNote != nil {
		h.quickNote.Show()
	}
}

// IsActive returns whether the hotkey is registered
func (h *TestManager) IsActive() bool {
	return h.isActive
} 