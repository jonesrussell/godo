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
		var _ Interface = (*Systray)(nil)
	})
}

func TestSystray(t *testing.T) {
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
		assert.True(t, svc.IsReady())
		assert.Equal(t, menu, svc.menu)
	})

	t.Run("SetIcon sets icon", func(t *testing.T) {
		svc := New(app, log)
		icon := &mockResource{}

		svc.SetIcon(icon)
		assert.Equal(t, icon, svc.icon)
		assert.False(t, svc.IsReady()) // Setting icon doesn't affect ready state
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
		assert.Nil(t, svc.menu)
		assert.Nil(t, svc.icon)

		// Test setup
		svc.Setup(menu)
		assert.True(t, svc.IsReady())
		assert.Equal(t, menu, svc.menu)

		// Test icon
		svc.SetIcon(icon)
		assert.Equal(t, icon, svc.icon)
		assert.True(t, svc.IsReady()) // Ready state should not be affected by icon
	})
}
