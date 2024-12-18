package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := initializeConfig()
	if err != nil {
		logger.Error("Failed to initialize config", "error", err)
		return
	}

	application, err := initializeApp(cfg)
	if err != nil {
		logger.Error("Failed to initialize app", "error", err)
		return
	}
	defer cleanup(application)

	// Load application icon
	iconBytes, err := assets.GetIcon()
	if err != nil {
		logger.Error("Failed to load application icon", "error", err)
	}

	fyneApp := fyneapp.New()
	fyneWin := fyneApp.NewWindow("Godo")
	fyneWin.Resize(fyne.NewSize(800, 600))
	fyneWin.CenterOnScreen()

	// Set up the content
	content := widget.NewLabel("Welcome to Godo")
	fyneWin.SetContent(content)

	// Set application icon if available
	if iconBytes != nil {
		icon := fyne.NewStaticResource("icon", iconBytes)
		fyneApp.SetIcon(icon)
	}

	// Register global shortcut
	if desk, ok := fyneApp.(desktop.App); ok {
		shortcut := &desktop.CustomShortcut{
			KeyName:  fyne.KeyG,
			Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt,
		}

		desk.SetSystemTrayMenu(fyne.NewMenu("Godo",
			fyne.NewMenuItem("Open", func() {
				fyneWin.Show()
				fyneWin.RequestFocus()
				fyneWin.CenterOnScreen()
			}),
			fyne.NewMenuItem("Quick Note", func() { showQuickNote(ctx, application) }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() { fyneApp.Quit() }),
		))

		fyneWin.Canvas().AddShortcut(shortcut, func(shortcut fyne.Shortcut) {
			logger.Debug("Global hotkey triggered")
			showQuickNote(ctx, application)
		})
	}

	// Run the application
	if err := runApplication(application); err != nil {
		logger.Error("Application error", "error", err)
		return
	}
}

func runApplication(application *app.App) error {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create an error channel
	errChan := make(chan error, 1)

	// Start the quick note feature in a goroutine
	go func() {
		if err := application.GetQuickNote().Show(ctx); err != nil {
			logger.Error("Quick note error", "error", err)
			errChan <- err
		}
	}()

	// Run the main UI in the main goroutine
	return application.Run(ctx)
}

func initializeConfig() (*config.Config, error) {
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		return nil, err
	}

	logConfig := common.LogConfig{
		Level:       cfg.Logging.Level,
		Output:      cfg.Logging.Output,
		ErrorOutput: cfg.Logging.ErrorOutput,
	}

	if _, err := logger.InitializeWithConfig(logConfig); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initializeApp(cfg *config.Config) (*app.App, error) {
	application, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	return application, nil
}

func cleanup(application *app.App) {
	logger.Info("Cleaning up application...")
	if err := application.Cleanup(); err != nil {
		logger.Error("Failed to cleanup", "error", err)
	}
}

func showQuickNote(ctx context.Context, application *app.App) {
	logger.Debug("Opening quick note window")

	quickNote, err := quicknote.New()
	if err != nil {
		logger.Error("Failed to create quick note UI", "error", err)
		return
	}

	inputChan := quickNote.GetInput()

	if err := quickNote.Show(ctx); err != nil {
		logger.Error("Failed to show quick note", "error", err)
		return
	}

	go func() {
		select {
		case input := <-inputChan:
			logger.Debug("Received quick note input, creating todo")
			todoService := application.GetTodoService()
			_, err := todoService.CreateTodo(ctx, input, "")
			if err != nil {
				logger.Error("Failed to create todo from quick note", "error", err)
			} else {
				logger.Debug("Successfully created todo from quick note")
			}
		case <-ctx.Done():
			logger.Debug("Quick note context cancelled")
			return
		}
	}()
}
