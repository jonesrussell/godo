// Package quicknote provides a quick note window for rapid task entry
package quicknote

import (
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
)

// Window represents a quick note window for rapid task entry
type Window struct {
	app    fyne.App
	window fyne.Window
	store  storage.TaskStore
	log    logger.Logger
	cfg    config.WindowConfig

	// UI components
	entry      *Entry
	addButton  *widget.Button
	clearBtn   *widget.Button
	statusText *widget.Label
}

// New creates a new quick note window
func New(
	app fyne.App,
	store storage.TaskStore,
	log logger.Logger,
	cfg config.WindowConfig,
) *Window {
	log.Debug("Creating new quick note window")
	w := &Window{
		app:    app,
		store:  store,
		log:    log,
		cfg:    cfg,
		window: app.NewWindow("Quick Note"),
	}

	w.setupUI()
	log.Debug("Quick note window created and UI setup completed")
	return w
}

// Initialize sets up the window with the given app and logger
func (w *Window) Initialize(app fyne.App, log logger.Logger) {
	w.app = app
	w.log = log
	w.log.Debug("Quick note window initialized")
}

// Show displays the quick note window
func (w *Window) Show() {
	w.log.Debug("Quick note window Show() called")
	w.log.Debug("Window state before Show", "window_nil", w.window == nil, "app_nil", w.app == nil)

	fyne.Do(func() {
		w.log.Debug("Inside fyne.Do - showing window")
		w.window.Show()
		w.log.Debug("Window Show() called")
		w.window.CenterOnScreen()
		w.log.Debug("Window centered")
		canvas := fyne.CurrentApp().Driver().CanvasForObject(w.entry)
		if canvas != nil {
			canvas.Focus(w.entry)
			w.log.Debug("Entry focused via canvas.Focus")
		} else {
			w.log.Warn("Could not get canvas for entry to focus")
		}
		w.log.Debug("Quick note window shown and focused")
	})
	w.log.Debug("Outside fyne.Do - Show() method completed")
}

// Hide hides the quick note window
func (w *Window) Hide() {
	fyne.Do(func() {
		w.window.Hide()
	})
}

// setupUI initializes the user interface
func (w *Window) setupUI() {
	w.log.Debug("Setting up quick note UI")

	// Create entry field
	w.entry = NewEntry()
	w.entry.SetPlaceHolder("Enter your task here...")
	w.entry.OnSubmitted = func(text string) {
		w.addTask()
	}

	// Create add button
	w.addButton = widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), w.addTask)

	// Create clear button
	w.clearBtn = widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), w.clearEntry)

	// Create status text
	w.statusText = widget.NewLabel("")
	w.statusText.Hide()

	// Create button container
	buttons := container.NewHBox(
		w.addButton,
		w.clearBtn,
		layout.NewSpacer(),
	)

	// Create main container
	content := container.NewVBox(
		w.entry,
		buttons,
		w.statusText,
	)

	w.window.SetContent(content)
	w.window.Resize(fyne.NewSize(400, 150))
	w.window.CenterOnScreen()

	// Set window properties
	w.window.SetCloseIntercept(func() {
		w.Hide()
	})

	w.log.Debug("Quick note UI setup completed")
}

// addTask adds a new task from the entry field
func (w *Window) addTask() {
	content := w.entry.Text
	if content == "" {
		return
	}

	// Create new task
	task := model.Task{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the task
	ctx := context.Background()
	if err := w.store.Add(ctx, &task); err != nil {
		w.log.Error("Failed to add task", "error", err)
		w.showStatus("Failed to add task")
		return
	}

	w.log.Info("Task added successfully", "task_id", task.ID)
	w.showStatus("Task added successfully!")
	w.clearEntry()
}

// clearEntry clears the entry field
func (w *Window) clearEntry() {
	w.entry.SetText("")
	w.statusText.Hide()
	canvas := fyne.CurrentApp().Driver().CanvasForObject(w.entry)
	if canvas != nil {
		canvas.Focus(w.entry)
		w.log.Debug("Entry focused via canvas.Focus (clearEntry)")
	} else {
		w.log.Warn("Could not get canvas for entry to focus (clearEntry)")
	}
}

// showStatus shows a status message
func (w *Window) showStatus(message string) {
	fyne.Do(func() {
		w.statusText.SetText(message)
		w.statusText.Show()
	})

	// Auto-hide after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		fyne.Do(func() {
			w.statusText.Hide()
		})
	}()
}
