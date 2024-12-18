package gui

import (
	"context"
	"os"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/logger"
)

type GUI struct {
	app        *app.App
	fyneApp    fyne.App
	mainWindow *MainWindow
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(application *app.App) *GUI {
	ctx, cancel := context.WithCancel(context.Background())
	return &GUI{
		app:     application,
		fyneApp: fyneapp.New(),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (g *GUI) Run() error {
	defer g.Cleanup()

	// Setup main window
	g.mainWindow = NewMainWindow(g)
	if err := g.mainWindow.Setup(); err != nil {
		return err
	}

	// Run the application
	g.fyneApp.Run()
	return nil
}

func (g *GUI) Cleanup() {
	logger.Info("GUI cleanup started")

	// Cancel context first
	g.cancel()

	// Close main window if it exists
	if g.mainWindow != nil {
		g.mainWindow.window.Close()
	}

	// Quit the Fyne app
	g.fyneApp.Quit()

	// Clean up the app
	if err := g.app.Cleanup(); err != nil {
		logger.Error("Failed to cleanup app", "error", err)
	}

	logger.Info("GUI cleanup completed")

	// Force exit if needed
	os.Exit(0)
}
