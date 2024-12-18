package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
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
	fyneWin := fyneApp.NewWindow("Godo Quick Note")

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
			fyne.NewMenuItem("Open", func() { fyneWin.Show() }),
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
	runApplication(ctx, cancel, application, fyneApp)
}

func runApplication(ctx context.Context, cancel context.CancelFunc, application *app.App, fyneApp fyne.App) {
	errChan := make(chan error, 1)

	// Start signal handling in a separate goroutine
	go handleSignals(ctx, cancel, errChan)

	// Run the application
	go func() {
		if err := application.Run(ctx); err != nil {
			logger.Error("Application error", "error", err)
			errChan <- fmt.Errorf("application error: %w", err)
			cancel()
		}
	}()

	// Run the Fyne app in the main thread
	fyneApp.Run()
}

func initializeConfig() (*config.Config, error) {
	env := os.Getenv("GODO_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.Load(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Convert config.LogConfig to common.LogConfig
	logConfig := common.LogConfig{
		Level:       cfg.Logging.Level,
		Output:      cfg.Logging.Output,
		ErrorOutput: cfg.Logging.ErrorOutput,
	}

	// Handle both return values from InitializeWithConfig
	if _, err := logger.InitializeWithConfig(logConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return cfg, nil
}

func initializeApp(cfg *config.Config) (*app.App, error) {
	application, err := app.InitializeAppWithConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize application: %w", err)
	}
	return application, nil
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, errChan <-chan error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer signal.Stop(sigChan)

	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGHUP:
				logger.Info("Received SIGHUP, reloading configuration...")
				// TODO: Implement config reload logic
				continue
			default:
				logger.Info("Initiating shutdown, received signal", "signal", sig)
				cancel()
				return
			}
		case err := <-errChan:
			logger.Error("Initiating shutdown due to error", "error", err)
			cancel()
			return
		case <-ctx.Done():
			logger.Info("Shutdown initiated by context cancellation")
			return
		}
	}
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
			// Get a pointer to TodoService before calling CreateTodo
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
