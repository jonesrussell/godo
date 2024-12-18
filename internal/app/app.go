// app.go
package app

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/quicknote"
	"github.com/jonesrussell/godo/internal/service"
	"github.com/jonesrussell/godo/internal/ui"
)

// App represents the main application
type App struct {
	config        *config.Config
	todoService   *service.TodoService
	hotkeyManager hotkey.HotkeyManager
	program       *tea.Program
	todoUI        *ui.TodoUI
	quickNote     quicknote.UI
}

// NewApp creates a new App instance with all dependencies
func NewApp(
	cfg *config.Config,
	todoService *service.TodoService,
	hotkeyManager hotkey.HotkeyManager,
	program *tea.Program,
	todoUI *ui.TodoUI,
	quickNote quicknote.UI,
) *App {
	return &App{
		config:        cfg,
		todoService:   todoService,
		hotkeyManager: hotkeyManager,
		program:       program,
		todoUI:        todoUI,
		quickNote:     quickNote,
	}
}

// GetTodoService returns the todo service instance
func (a *App) GetTodoService() *service.TodoService {
	return a.todoService
}

// GetHotkeyManager returns the hotkey manager instance
func (a *App) GetHotkeyManager() hotkey.HotkeyManager {
	return a.hotkeyManager
}

// GetProgram returns the Bubble Tea program instance
func (a *App) GetProgram() *tea.Program {
	return a.program
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	logger.Info("Starting application services...")

	// Start the hotkey manager
	if err := a.hotkeyManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start hotkey manager: %w", err)
	}

	// Register the quick note hotkey
	if a.config.Hotkeys.QuickNote == nil {
		return fmt.Errorf("quick note hotkey configuration is missing")
	}

	if err := a.hotkeyManager.RegisterHotkey(*a.config.Hotkeys.QuickNote); err != nil {
		return fmt.Errorf("failed to register quick note hotkey: %w", err)
	}

	logger.Info("Hotkey registered successfully: %v", a.config.Hotkeys.QuickNote)

	if err := a.initializeServices(ctx); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Start background service to handle hotkey events
	go func() {
		hotkeyEvents := a.hotkeyManager.GetEventChannel()
		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopping hotkey listener...")
				return
			case <-hotkeyEvents:
				logger.Info("Hotkey triggered - showing quick note")
				if err := a.handleQuickNote(ctx); err != nil {
					logger.Error("Failed to handle quick note", "error", err)
				}
			}
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return ctx.Err()
}

func (a *App) handleQuickNote(ctx context.Context) error {
	logger.Info("Quick note triggered")

	if err := a.quickNote.Show(ctx); err != nil {
		return fmt.Errorf("failed to show quick note: %w", err)
	}

	return nil
}

func (a *App) initializeServices(ctx context.Context) error {
	logger.Info("Initializing services...")

	// Verify database connection
	testTodo, err := a.todoService.CreateTodo(ctx, "test", "Testing service initialization")
	if err != nil {
		logger.Error("Failed to verify database connection: %v", err)
		return fmt.Errorf("failed to verify database connection: %w", err)
	}

	// Clean up test todo
	if err := a.todoService.DeleteTodo(ctx, testTodo.ID); err != nil {
		logger.Error("Failed to cleanup test todo: %v", err)
		return fmt.Errorf("failed to cleanup test todo: %w", err)
	}

	logger.Info("Services initialized successfully")
	return nil
}

// Cleanup performs any necessary cleanup before shutdown
func (a *App) Cleanup() error {
	logger.Info("Cleaning up application resources")
	// Add any cleanup logic here
	return nil
}
