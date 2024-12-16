package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2/app"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/getlantern/systray"
	"github.com/jonesrussell/godo/internal/di"
	"github.com/jonesrussell/godo/internal/icon"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

var (
	fullUI = flag.Bool("ui", false, "Launch full todo management interface")
)

// setupSignalHandler creates a signal handler
func setupSignalHandler(parentCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(parentCtx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer cancel()
		select {
		case sig := <-sigChan:
			logger.Info("Received signal: %v", sig)
			cancel()
		case <-parentCtx.Done():
			// Parent context was cancelled
		}
	}()

	return ctx
}

// onReady is called when systray is ready
func onReady(ctx context.Context, app *di.App, cancel context.CancelFunc) func() {
	return func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Godo")
		systray.SetTooltip("Quick Todo Manager")

		systray.AddSeparator()
		mOpen := systray.AddMenuItem("Open Manager", "Open todo manager")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit application")

		// Start background service
		go func() {
			if err := app.GetHotkeyManager().Start(ctx); err != nil {
				logger.Error("Failed to start hotkey manager: %v", err)
				return
			}

			hotkeyEvents := app.GetHotkeyManager().GetEventChannel()
			logger.Info("Listening for hotkey events (Ctrl+Alt+G)...")

			for {
				select {
				case <-ctx.Done():
					if err := app.GetHotkeyManager().Cleanup(); err != nil {
						logger.Error("Error cleaning up hotkey: %v", err)
					}
					return
				case <-hotkeyEvents:
					logger.Info("Hotkey triggered - showing quick note")
					quickNote := ui.NewQuickNote(app.GetTodoService(), app.GetFyneApp())
					quickNote.Show()
				}
			}
		}()

		// Handle menu items
		for {
			select {
			case <-mQuit.ClickedCh:
				cancel()
				systray.Quit()
				return
			case <-mOpen.ClickedCh:
				showFullUI(app.GetTodoService())
			case <-ctx.Done():
				systray.Quit()
				return
			}
		}
	}
}

// onExit is called when systray is quitting
func onExit() {
	logger.Info("Cleaning up...")
	os.Exit(0)
}

// showFullUI displays the full todo management interface
func showFullUI(service *service.TodoService) {
	p := tea.NewProgram(
		ui.New(service),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		logger.Error("UI error: %v", err)
	}
}

func main() {
	cleanup := logger.Initialize()
	defer cleanup()

	flag.Parse()

	logger.Info("Starting Godo application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCtx := setupSignalHandler(ctx)

	diApp, err := di.InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize application: %v", err)
	}

	fyneApp := app.New()
	diApp.SetFyneApp(fyneApp)

	// Create main window but don't show it
	mainWindow := fyneApp.NewWindow("Godo")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	if *fullUI {
		showFullUI(diApp.GetTodoService())
	} else {
		// Start systray first
		go systray.Run(onReady(sigCtx, diApp, cancel), onExit)

		// Run Fyne app in the main thread
		go func() {
			<-sigCtx.Done()
			logger.Info("Starting graceful shutdown...")
			fyneApp.Quit()
		}()

		// This will block until the app quits
		fyneApp.Run()
	}
}
