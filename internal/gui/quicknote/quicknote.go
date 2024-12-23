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

type customEntry struct {
	widget.Entry
	onCtrlEnter func()
	onEscape    func()
	logger      logger.Logger
}

func newCustomEntry(log logger.Logger) *customEntry {
	entry := &customEntry{
		logger: log,
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *customEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if cs, ok := shortcut.(*desktop.CustomShortcut); ok {
		if cs.KeyName == fyne.KeyReturn && cs.Modifier == fyne.KeyModifierControl {
			e.logger.Debug("Ctrl+Enter shortcut triggered")
			e.onCtrlEnter()
			return
		}
	}
	e.Entry.TypedShortcut(shortcut)
}

func (e *customEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		e.logger.Debug("Escape key pressed")
		e.onEscape()
		return
	}
	e.Entry.TypedKey(key)
}

type QuickNote struct {
	window     fyne.Window
	mainWindow fyne.Window
	store      storage.Store
	input      *customEntry
	form       dialog.Dialog
	logger     logger.Logger
}

// New creates a new QuickNote instance
func New(app fyne.App, mainWindow fyne.Window, store storage.Store, log logger.Logger) *QuickNote {
	qn := &QuickNote{
		mainWindow: mainWindow,
		window:     app.NewWindow("Quick Note"),
		store:      store,
		logger:     log,
	}

	// Set a reasonable default size for the quick note window
	qn.window.Resize(fyne.NewSize(400, 200))

	qn.input = newCustomEntry(log)
	qn.setupUI()

	return qn
}

func (qn *QuickNote) setupUI() {
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

	qn.form.Resize(fyne.NewSize(380, 180))
}

func (qn *QuickNote) setupShortcuts() {
	qn.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		qn.logger.Debug("Window Ctrl+Enter shortcut triggered")
		qn.handleSave()
	})
}

func (qn *QuickNote) handleSave() {
	if qn.input.Text != "" {
		todo := model.NewTodo(qn.input.Text)
		if err := qn.store.Add(todo); err != nil {
			qn.logger.Error("Failed to save todo", "error", err)
			dialog.ShowError(err, qn.window)
		} else {
			qn.logger.Debug("Saved note as todo", "id", todo.ID, "content", todo.Content)
		}
	}
	qn.handleCancel()
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
	qn.mainWindow.Hide()
	qn.window.Show()
	qn.window.CenterOnScreen()
	qn.form.Show()
	qn.window.Canvas().Focus(qn.input)
}

// Hide hides the quick note dialog
func (qn *QuickNote) Hide() {
	qn.input.SetText("")
	qn.window.Hide()
	qn.logger.Debug("Quick note hidden")
}
