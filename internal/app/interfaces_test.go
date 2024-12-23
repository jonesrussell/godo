package app_test

import (
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestMockUI(t *testing.T) {
	test.NewApp()

	tests := []struct {
		name     string
		setup    func() *app.MockUI
		action   func(*app.MockUI)
		validate func(*testing.T, *app.MockUI)
	}{
		{
			name: "Show sets show to true",
			setup: func() *app.MockUI {
				return &app.MockUI{}
			},
			action: func(m *app.MockUI) {
				m.Show()
			},
			validate: func(_ *testing.T, m *app.MockUI) {
				assert.True(t, m.IsShown())
			},
		},
		{
			name: "Hide sets show to false",
			setup: func() *app.MockUI {
				m := &app.MockUI{}
				m.Show() // Start shown
				return m
			},
			action: func(m *app.MockUI) {
				m.Hide()
			},
			validate: func(t *testing.T, m *app.MockUI) {
				assert.False(t, m.IsShown())
			},
		},
		{
			name: "SetContent stores content",
			setup: func() *app.MockUI {
				return &app.MockUI{}
			},
			action: func(m *app.MockUI) {
				content := canvas.NewText("Test", test.Theme().Color(theme.ColorNameForeground, theme.VariantLight))
				m.SetContent(content)
			},
			validate: func(t *testing.T, m *app.MockUI) {
				assert.NotNil(t, m.Content())
				assert.IsType(t, &canvas.Text{}, m.Content())
			},
		},
		{
			name: "Resize stores size",
			setup: func() *app.MockUI {
				return &app.MockUI{}
			},
			action: func(m *app.MockUI) {
				size := fyne.NewSize(800, 600)
				m.Resize(size)
			},
			validate: func(t *testing.T, m *app.MockUI) {
				expected := fyne.NewSize(800, 600)
				assert.Equal(t, expected, m.Size())
			},
		},
		{
			name: "CenterOnScreen is no-op",
			setup: func() *app.MockUI {
				return &app.MockUI{}
			},
			action: func(m *app.MockUI) {
				m.CenterOnScreen()
			},
			validate: func(_ *testing.T, _ *app.MockUI) {
				// No validation needed for no-op
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setup()
			tt.action(mock)
			tt.validate(t, mock)
		})
	}
}
