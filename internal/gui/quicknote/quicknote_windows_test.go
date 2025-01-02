//go:build windows

package quicknote

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestQuickNoteHotkey(t *testing.T) {
	store := storage.NewMockStore()
	log := logger.NewMockTestLogger(t)
	app := test.NewApp()
	cfg := config.WindowConfig{
		Width:       200,
		Height:      100,
		StartHidden: false,
	}
	quickNote := New(app, store, log, cfg)
	assert.NotNil(t, quickNote)
}
