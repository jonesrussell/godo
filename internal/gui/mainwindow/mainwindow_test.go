package systray

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
)

type mockResource struct {
	fyne.Resource
}

func TestInterface(t *testing.T) {
	t.Run("implementation satisfies interface", func(_ *testing.T) {
		var _ Interface = (*Service)(nil)
	})
}

func TestService(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
	})
	assert.NoError(t, err)

	t.Run("Setup sets menu", func(t *testing.T) {
		svc := New(app, log)
		menu := &fyne.Menu{
			Label: "Test Menu",
			Items: []*fyne.MenuItem{
				{Label: "Test Item"},
			},
		}

		svc.Setup(menu)
		// We can only test the behavior since menu is private
		assert.True(t, svc.IsReady())
	})

	t.Run("SetIcon sets icon", func(t *testing.T) {
		svc := New(app, log)
		icon := &mockResource{}

		svc.SetIcon(icon)
		// We can only test the behavior since icon is private
		assert.True(t, svc.IsReady())
	})

	t.Run("IsReady returns ready state", func(t *testing.T) {
		svc := New(app, log)
		assert.False(t, svc.IsReady())

		svc.Setup(&fyne.Menu{})
		assert.True(t, svc.IsReady())
	})
}

func TestIntegration(t *testing.T) {
	t.Run("full lifecycle", func(t *testing.T) {
		app := test.NewApp()
		defer app.Quit()

		log, err := logger.New(&common.LogConfig{
			Level:   "debug",
			Console: true,
		})
		assert.NoError(t, err)

		svc := New(app, log)
		menu := &fyne.Menu{
			Label: "Test Menu",
			Items: []*fyne.MenuItem{
				{Label: "Test Item"},
			},
		}
		icon := &mockResource{}

		// Test initial state
		assert.False(t, svc.IsReady())

		// Test setup
		svc.Setup(menu)
		assert.True(t, svc.IsReady())

		// Test icon
		svc.SetIcon(icon)
		assert.True(t, svc.IsReady())
	})
}
