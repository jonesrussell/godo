package quicknote

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
)

// Config holds all dependencies for QuickNote
type Config struct {
	App        fyne.App
	MainWindow fyne.Window
	Store      storage.Store
	Logger     logger.Logger
}

// QuickNote manages the quick note window and functionality
type QuickNote struct {
	config     Config
	window     fyne.Window
	input      *customEntry
	form       dialog.Dialog
	dimensions struct {
		window fyne.Size
		form   fyne.Size
	}
}

// New creates a new QuickNote instance
func New(cfg Config) *QuickNote {
	qn := &QuickNote{
		config: cfg,
		window: cfg.App.NewWindow("Quick Note"),
		dimensions: struct {
			window fyne.Size
			form   fyne.Size
		}{
			window: fyne.NewSize(400, 200),
			form:   fyne.NewSize(380, 180),
		},
	}

	qn.input = newCustomEntry(cfg.Logger)
	qn.setupUI()

	return qn
}

func (qn *QuickNote) setupUI() {
	qn.window.Resize(qn.dimensions.window)
	qn.setupInput()
	qn.setupForm()
	qn.setupShortcuts()
}

func (qn *QuickNote) setupInput() {
	qn.input.SetPlaceHolder("Enter your note here...")
	qn.input.onCtrlEnter = qn.handleSave
	qn.input.onEscape = qn.handleCancel
}

func (qn *QuickNote) setupForm() {
	qn.form = dialog.NewForm("Quick Note", "Save", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Note", qn.input),
		},
		qn.handleFormSubmit,
		qn.window)

	qn.form.Resize(qn.dimensions.form)
}

func (qn *QuickNote) setupShortcuts() {
	qn.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		qn.config.Logger.Debug("Window Ctrl+Enter shortcut triggered")
		qn.handleSave()
	})
}

func (qn *QuickNote) handleSave() {
	if qn.input.Text != "" {
		if err := qn.saveTodo(qn.input.Text); err != nil {
			qn.config.Logger.Error("Failed to save todo", "error", err)
			dialog.ShowError(err, qn.window)
			return
		}
		qn.config.Logger.Debug("Saved note as todo", "content", qn.input.Text)
	}
	qn.handleCancel()
}

func (qn *QuickNote) saveTodo(content string) error {
	todo := model.NewTodo(content)
	return qn.config.Store.Add(todo)
}

func (qn *QuickNote) handleCancel() {
	qn.input.SetText("")
	qn.Hide()
}

func (qn *QuickNote) handleFormSubmit(save bool) {
	if save {
		qn.handleSave()
	} else {
		qn.handleCancel()
	}
}

// Show displays the quick note dialog
func (qn *QuickNote) Show() {
	qn.config.MainWindow.Hide()
	qn.window.Show()
	qn.window.CenterOnScreen()
	qn.form.Show()
	qn.window.Canvas().Focus(qn.input)
}

// Hide hides the quick note dialog
func (qn *QuickNote) Hide() {
	qn.input.SetText("")
	qn.window.Hide()
	qn.config.Logger.Debug("Quick note hidden")
}

// GetWindow returns the underlying window for testing
func (qn *QuickNote) GetWindow() fyne.Window {
	return qn.window
}

// GetInput returns the input field for testing
func (qn *QuickNote) GetInput() *customEntry {
	return qn.input
}
