// Package app implements the main application logic for Godo.
package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/api"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// App represents the main application
type App struct {
	store    storage.Store
	window   gui.MainWindowManager
	logger   logger.Logger
	stopOnce sync.Once
	stopChan chan struct{}
	notes    map[string]storage.Note
	list     *widget.List
}

// SetupUI implements the UIManager interface
func (a *App) SetupUI() error {
	return a.initWindow()
}

// Run implements the ApplicationService interface
func (a *App) Run() {
	a.window.Show()
}

// Cleanup implements the ApplicationService interface
func (a *App) Cleanup() {
	_ = a.Stop()
}

// Store returns the storage.Store instance
func (a *App) Store() storage.Store {
	return a.store
}

// Logger returns the logger instance
func (a *App) Logger() logger.Logger {
	return a.logger
}

// Params contains the parameters for creating a new App instance
type Params struct {
	Store     storage.Store
	Window    gui.MainWindowManager
	Logger    logger.Logger
	Hotkey    hotkey.Manager
	APIServer *api.Server
	APIRunner *api.Runner
}

// New creates a new App instance
func New(params Params) (*App, error) {
	if params.Store == nil {
		return nil, fmt.Errorf("store is required")
	}
	if params.Window == nil {
		return nil, fmt.Errorf("window is required")
	}
	if params.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &App{
		store:    params.Store,
		window:   params.Window,
		logger:   params.Logger,
		stopChan: make(chan struct{}),
		notes:    make(map[string]storage.Note),
	}, nil
}

// Start starts the application
func (a *App) Start() error {
	// Initialize the window
	if err := a.initWindow(); err != nil {
		return fmt.Errorf("failed to initialize window: %w", err)
	}

	// Show the window
	a.window.Show()

	return nil
}

// Stop stops the application
func (a *App) Stop() error {
	a.stopOnce.Do(func() {
		close(a.stopChan)
	})

	if err := a.store.Close(); err != nil {
		return fmt.Errorf("failed to close store: %w", err)
	}

	return nil
}

// initWindow initializes the main window
func (a *App) initWindow() error {
	// Load initial notes
	notes, err := a.store.List(context.Background())
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	// Store notes in memory
	for _, note := range notes {
		a.notes[note.ID] = note
	}

	// Create list widget
	a.list = widget.NewList(
		func() int { return len(a.notes) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			note := a.getNoteByIndex(id)
			if note == nil {
				return
			}

			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)

			check.Checked = note.Completed
			check.OnChanged = func(checked bool) {
				note.Completed = checked
				note.UpdatedAt = time.Now().Unix()
				a.UpdateNote(context.Background(), *note)
			}

			label.Text = note.Content
			label.Refresh()
		},
	)

	// Create input field for new notes
	input := widget.NewEntry()
	input.SetPlaceHolder("Add a new note...")
	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}

		now := time.Now().Unix()
		note := storage.Note{
			ID:        fmt.Sprintf("%d", now),
			Content:   text,
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := a.AddNote(context.Background(), note); err != nil {
			// Handle error (in a real app, show error to user)
			return
		}

		input.SetText("")
	}

	// Create main content
	content := container.NewBorder(input, nil, nil, nil, a.list)
	a.window.SetContent(content)

	return nil
}

// getNoteByIndex returns a note by its list index
func (a *App) getNoteByIndex(index int) *storage.Note {
	i := 0
	for _, note := range a.notes {
		if i == index {
			return &note
		}
		i++
	}
	return nil
}

// AddNote adds a new note
func (a *App) AddNote(ctx context.Context, note storage.Note) error {
	if err := a.store.Add(ctx, note); err != nil {
		return fmt.Errorf("failed to add note: %w", err)
	}

	// Update in-memory notes
	a.notes[note.ID] = note
	a.list.Refresh()

	return nil
}

// UpdateNote updates an existing note
func (a *App) UpdateNote(ctx context.Context, note storage.Note) error {
	if err := a.store.Update(ctx, note); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	// Update in-memory notes
	a.notes[note.ID] = note
	a.list.Refresh()

	return nil
}

// DeleteNote deletes a note by ID
func (a *App) DeleteNote(ctx context.Context, id string) error {
	if err := a.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Update in-memory notes
	delete(a.notes, id)
	a.list.Refresh()

	return nil
}
