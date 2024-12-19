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
}

func newCustomEntry() *customEntry {
	entry := &customEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *customEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if cs, ok := shortcut.(*desktop.CustomShortcut); ok {
		if cs.KeyName == fyne.KeyReturn && cs.Modifier == fyne.KeyModifierControl {
			logger.Debug("Ctrl+Enter shortcut triggered")
			e.onCtrlEnter()
			return
		}
	}
	e.Entry.TypedShortcut(shortcut)
}

func (e *customEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyEscape {
		logger.Debug("Escape key pressed")
		e.onEscape()
		return
	}
	e.Entry.TypedKey(key)
}

// QuickNote represents a quick note dialog
type QuickNote struct {
	window fyne.Window
	store  storage.Store
	input  *customEntry
	form   dialog.Dialog
}

// New creates a new QuickNote instance
func New(window fyne.Window, store storage.Store) *QuickNote {
	qn := &QuickNote{
		window: window,
		store:  store,
		input:  newCustomEntry(),
	}

	qn.setupInput()
	qn.setupForm()
	qn.setupShortcuts()

	return qn
}

func (qn *QuickNote) setupInput() {
	qn.input.SetPlaceHolder("Enter your note here...")
}

func (qn *QuickNote) setupForm() {
	qn.form = dialog.NewForm("Quick Note", "Save", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Note", qn.input),
		},
		qn.handleFormSubmit,
		qn.window)

	qn.form.Resize(fyne.NewSize(400, 200))
}

func (qn *QuickNote) setupShortcuts() {
	qn.input.onCtrlEnter = func() {
		if qn.input.Text != "" {
			todo := model.NewTodo(qn.input.Text)
			if err := qn.store.Add(todo); err != nil {
				logger.Error("Failed to save todo", "error", err)
				dialog.ShowError(err, qn.window)
			} else {
				logger.Debug("Saved note as todo", "id", todo.ID, "content", todo.Content)
			}
		}
		qn.input.SetText("")
		qn.window.Hide()
		logger.Debug("Quick note saved and window hidden")
	}

	qn.input.onEscape = func() {
		qn.input.SetText("")
		qn.window.Hide()
		logger.Debug("Quick note cancelled and window hidden")
	}

	// Register Ctrl+Enter with the window
	qn.window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyReturn,
		Modifier: fyne.KeyModifierControl,
	}, func(shortcut fyne.Shortcut) {
		logger.Debug("Window Ctrl+Enter shortcut triggered")
		qn.input.onCtrlEnter()
	})
}

func (qn *QuickNote) handleFormSubmit(save bool) {
	if save {
		qn.saveTodo()
	} else {
		qn.cancel()
	}
}

func (qn *QuickNote) saveTodo() {
	if qn.input.Text != "" {
		todo := model.NewTodo(qn.input.Text)
		if err := qn.store.Add(todo); err != nil {
			logger.Error("Failed to save todo", "error", err)
			dialog.ShowError(err, qn.window)
		} else {
			logger.Debug("Saved note as todo", "id", todo.ID, "content", todo.Content)
		}
	}
	qn.cancel()
}

func (qn *QuickNote) cancel() {
	qn.input.SetText("")
	qn.window.Hide()
}

// Show displays the quick note dialog
func (qn *QuickNote) Show() {
	qn.window.Show()
	qn.window.CenterOnScreen()
	qn.form.Show()
	qn.window.Canvas().Focus(qn.input)
}

// Hide hides the quick note dialog
func (qn *QuickNote) Hide() {
	qn.input.SetText("")
	qn.window.Hide()
	logger.Debug("Quick note hidden")
}
