package container

import (
	"fmt"

	godoapp "github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// Container holds all application dependencies
type Container struct {
	App    godoapp.ApplicationService
	Logger logger.Logger
	Store  types.Store
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
	godoApp := app.(*godoapp.App)
	return &Container{
		App:    app,
		Logger: godoApp.Logger(),
		Store:  godoApp.Store(),
	}, nil
}
