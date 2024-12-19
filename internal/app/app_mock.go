package app

import (
	"fyne.io/fyne/v2"
)

type MockUI struct {
	content fyne.CanvasObject
	size    fyne.Size
	show    bool
}

func (m *MockUI) Show() {
	m.show = true
}

func (m *MockUI) Hide() {
	m.show = false
}

func (m *MockUI) SetContent(content fyne.CanvasObject) {
	m.content = content
}

func (m *MockUI) Resize(size fyne.Size) {
	m.size = size
}

func (m *MockUI) CenterOnScreen() {
	// No-op for mock
}

// Getter methods for testing
func (m *MockUI) IsShown() bool {
	return m.show
}

func (m *MockUI) Content() fyne.CanvasObject {
	return m.content
}

func (m *MockUI) Size() fyne.Size {
	return m.size
}
