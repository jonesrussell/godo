package container

import (
	"fmt"

	"github.com/jonesrussell/godo/internal/application/core"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// Container holds all application dependencies
type Container struct {
	App    core.Application
	Logger logger.Logger
	Store  storage.TaskStore
}

// New creates a new container instance
func New() (*Container, error) {
	app, cleanup, err := InitializeApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application: %w", err)
	}
	defer func() {
		if err != nil {
			cleanup()
		}
	}()

	// Get logger and store from the app
	godoApp, ok := app.(*core.App)
	if !ok {
		cleanup()
		return nil, fmt.Errorf("failed to cast app to *godoapp.App: unexpected type %T", app)
	}
	return &Container{
		App:    app,
		Logger: godoApp.Logger(),
		Store:  godoApp.Store(),
	}, nil
}
