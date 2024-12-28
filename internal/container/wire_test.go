//go:build wireinject && windows

package container

import (
	"testing"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestInitializeTestApp(t *testing.T) {
	// First verify we can create all the individual dependencies
	log, cleanup, err := ProvideLogger(&LoggerOptions{
		Level:       ProvideLogLevel(),
		Output:      ProvideLogOutputPaths(),
		ErrorOutput: ProvideErrorOutputPaths(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, log)
	defer cleanup()

	store := ProvideMockStore()
	assert.NotNil(t, store)

	mainWin := ProvideMockMainWindow()
	assert.NotNil(t, mainWin)

	quickNote := ProvideMockQuickNote()
	assert.NotNil(t, quickNote)

	hotkey := ProvideMockHotkey()
	assert.NotNil(t, hotkey)

	httpConfig := ProvideHTTPConfig(&HTTPOptions{
		Port:              ProvideHTTPPort(),
		ReadTimeout:       ProvideReadTimeout(),
		WriteTimeout:      ProvideWriteTimeout(),
		ReadHeaderTimeout: ProvideHeaderTimeout(),
		IdleTimeout:       ProvideIdleTimeout(),
	})
	assert.NotNil(t, httpConfig)

	// Now test the full app initialization
	testApp := &app.TestApp{
		Logger:     log,
		Store:      store,
		MainWindow: mainWin,
		QuickNote:  quickNote,
		Hotkey:     hotkey,
		HTTPConfig: httpConfig,
		Name:       ProvideAppName(),
		Version:    ProvideAppVersion(),
		ID:         ProvideAppID(),
	}
	assert.NotNil(t, testApp)
	assert.NotNil(t, testApp.Logger)
	assert.NotNil(t, testApp.Store)
	assert.NotNil(t, testApp.MainWindow)
	assert.NotNil(t, testApp.QuickNote)
	assert.NotNil(t, testApp.Hotkey)
	assert.NotNil(t, testApp.HTTPConfig)
	assert.Equal(t, "Godo", testApp.Name.String())
	assert.Equal(t, "1.0.0", testApp.Version.String())
	assert.Equal(t, "com.jonesrussell.godo", testApp.ID.String())
}

func TestProvideAppName(t *testing.T) {
	name := ProvideAppName()
	assert.Equal(t, "Godo", name.String())
}

func TestProvideAppVersion(t *testing.T) {
	version := ProvideAppVersion()
	assert.NotEmpty(t, version.String(), "Version should not be empty")
	assert.Regexp(t, `^\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$`, version.String(), "Version should follow semantic versioning")
}

func TestProvideAppID(t *testing.T) {
	id := ProvideAppID()
	assert.Equal(t, "com.jonesrussell.godo", id.String())
}

func TestProvideDatabasePath(t *testing.T) {
	path := ProvideDatabasePath()
	assert.Equal(t, "godo.db", path.String())
}

func TestProvideLogLevel(t *testing.T) {
	level := ProvideLogLevel()
	assert.Equal(t, "info", level.String())
}
