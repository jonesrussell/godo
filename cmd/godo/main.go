package main

import (
	"context"
	"os"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/assets"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/marcsauter/single"
)

func main() {
	// Initialize logger first
	if _, err := logger.Initialize(); err != nil {
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Create single instance lock
	s := single.New("godo")
	if err := s.CheckLock(); err != nil {
		if err == single.ErrAlreadyRunning {
			logger.Info("Godo is already running. Look for the icon in your system tray.")
			return
		}
		logger.Error("Failed to check application lock", "error", err)
		return
	}
	defer func() {
		if err := s.TryUnlock(); err != nil {
			logger.Error("Failed to unlock single instance", "error", err)
		}
	}()

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

	// Register system tray
	if desk, ok := fyneApp.(desktop.App); ok {
		menu := fyne.NewMenu("Godo",
			fyne.NewMenuItem("Open", func() {
				logger.Debug("Opening main window")
				fyneWin.Show()
				fyneWin.RequestFocus()
				fyneWin.CenterOnScreen()
			}),
			fyne.NewMenuItem("Quick Note", func() {
				logger.Debug("Opening quick note from tray")
				showQuickNote(ctx, application)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				logger.Info("Quitting application")
				fyneApp.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(fyneApp.Icon())
	} else {
		logger.Warn("System tray not supported on this platform")
	}

	// Run the application
	fyneApp.Run()
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

	// Create a new context with cancellation for this quick note instance
	qnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	quickNote := application.GetQuickNote() // Use the existing QuickNote instance
	if quickNote == nil {
		logger.Error("Failed to get quick note instance")
		return
	}

	// Show the quick note window in a goroutine
	go func() {
		if err := quickNote.Show(qnCtx); err != nil {
			logger.Error("Failed to show quick note", "error", err)
			return
		}
	}()

	// Handle input in a separate goroutine
	go func() {
		select {
		case input := <-quickNote.GetInput():
			logger.Debug("Received quick note input", "input", input)
			todoService := application.GetTodoService()
			_, err := todoService.CreateTodo(qnCtx, input, "")
			if err != nil {
				logger.Error("Failed to create todo from quick note", "error", err)
			} else {
				logger.Debug("Successfully created todo from quick note")
			}
		case <-qnCtx.Done():
			logger.Debug("Quick note context cancelled")
			return
		}
	}()
}
