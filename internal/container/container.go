package container

import (
	"fyne.io/fyne/v2/app"
	godoapp "github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
	"go.uber.org/zap"
)

// Container holds all application dependencies
type Container struct {
	App    *godoapp.App
	Logger logger.Logger
	Store  storage.Store
}

// Initialize creates a new container with all dependencies
func Initialize(zapLogger *zap.SugaredLogger) (*Container, error) {
	// Create logger adapter
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}
	log, err := logger.New(logConfig)
	if err != nil {
		return nil, err
	}

	// Create SQLite store
	store, err := sqlite.New("godo.db", zapLogger.Desugar())
	if err != nil {
		return nil, err
	}

	// Create Fyne app
	fyneApp := app.NewWithID("com.jonesrussell.godo")

	// Create main window
	var mainWin gui.MainWindow = mainwindow.New(store, log)

	// Create quick note
	quickNote := quicknote.New(store, log)

	// Create hotkey manager
	hotkeys := hotkey.New(quickNote)

	// Create app
	godoApp := godoapp.New(log, fyneApp, store, mainWin, quickNote, hotkeys)

	return &Container{
		App:    godoApp,
		Logger: log,
		Store:  store,
	}, nil
}
