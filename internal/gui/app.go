package gui

import (
	"context"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"github.com/jonesrussell/godo/internal/app"
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
	g.cancel()
	g.fyneApp.Quit()
}
