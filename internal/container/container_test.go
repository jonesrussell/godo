//go:build !docker && wireinject && windows

package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testApp, cleanup, err := InitializeTestApp()
	require.NoError(t, err)
	require.NotNil(t, testApp)
	defer cleanup()

	// Verify all dependencies are properly initialized
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
