// Package windowmanager provides centralized window management for the application
package windowmanager

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

// FocusManager manages focus for windows and their components
type FocusManager struct {
	window     fyne.Window
	log        logger.Logger
	focusQueue []fyne.Focusable
	queueMu    sync.RWMutex
}

// NewFocusManager creates a new focus manager
func NewFocusManager(window fyne.Window, log logger.Logger) *FocusManager {
	return &FocusManager{
		window:     window,
		log:        log,
		focusQueue: make([]fyne.Focusable, 0),
	}
}

// RequestFocus requests focus for the window with multiple fallback strategies
func (fm *FocusManager) RequestFocus() {
	fyne.Do(func() {
		if fm.window == nil {
			fm.log.Warn("Cannot request focus: window is nil")
			return
		}

		// Strategy 1: Try to focus the window first
		fm.window.RequestFocus()

		// Strategy 2: Try to focus the first focusable element
		if err := fm.focusFirstElement(); err != nil {
			fm.log.Debug("Failed to focus first element, trying delayed focus", "error", err)
			// Strategy 3: Try delayed focus
			go fm.delayedFocus()
		}

		fm.log.Debug("Focus request completed")
	})
}

// FocusElement focuses a specific element with fallback strategies
func (fm *FocusManager) FocusElement(element fyne.Focusable) error {
	if element == nil {
		return fmt.Errorf("element is nil")
	}

	fyne.Do(func() {
		// Strategy 1: Direct focus
		if err := fm.focusDirect(element); err != nil {
			fm.log.Debug("Direct focus failed, trying canvas focus", "error", err)

			// Strategy 2: Canvas focus
			if err := fm.focusViaCanvas(element); err != nil {
				fm.log.Debug("Canvas focus failed, trying delayed focus", "error", err)

				// Strategy 3: Delayed focus
				go fm.delayedFocusElement(element)
			}
		}
	})

	return nil
}

// AddToFocusQueue adds an element to the focus queue
func (fm *FocusManager) AddToFocusQueue(element fyne.Focusable) {
	fm.queueMu.Lock()
	defer fm.queueMu.Unlock()

	fm.focusQueue = append(fm.focusQueue, element)
}

// ClearFocusQueue clears the focus queue
func (fm *FocusManager) ClearFocusQueue() {
	fm.queueMu.Lock()
	defer fm.queueMu.Unlock()

	fm.focusQueue = make([]fyne.Focusable, 0)
}

// HandleTabKey handles tab key navigation
func (fm *FocusManager) HandleTabKey(ke *fyne.KeyEvent) {
	fm.queueMu.RLock()
	queue := fm.focusQueue
	fm.queueMu.RUnlock()

	if len(queue) == 0 {
		return
	}

	// Find current focused element by comparing pointer addresses
	var currentIndex int = -1
	for i, element := range queue {
		if fyne.CurrentApp().Driver().CanvasForObject(element.(fyne.CanvasObject)).Focused() == element {
			currentIndex = i
			break
		}
	}

	// Default to first element if none focused
	if currentIndex == -1 {
		currentIndex = 0
	}

	// Calculate next index (Tab = forward, Shift+Tab = backward)
	// Fyne KeyEvent does not have Modifiers, so always go forward
	nextIndex := (currentIndex + 1) % len(queue)

	// Focus next element
	fm.FocusElement(queue[nextIndex])
}

// focusFirstElement tries to focus the first focusable element
func (fm *FocusManager) focusFirstElement() error {
	fm.queueMu.RLock()
	queue := fm.focusQueue
	fm.queueMu.RUnlock()

	if len(queue) == 0 {
		return fmt.Errorf("no focusable elements in queue")
	}

	return fm.focusDirect(queue[0])
}

// focusDirect attempts direct focus on an element
func (fm *FocusManager) focusDirect(element fyne.Focusable) error {
	if element == nil {
		return fmt.Errorf("element is nil")
	}

	// Try to focus the element
	element.FocusGained()
	return nil
}

// focusViaCanvas attempts to focus via the canvas
func (fm *FocusManager) focusViaCanvas(element fyne.Focusable) error {
	if fm.window == nil {
		return fmt.Errorf("window is nil")
	}

	canvas := fm.window.Canvas()
	if canvas == nil {
		return fmt.Errorf("canvas is nil")
	}

	// Try to focus via canvas
	canvas.Focus(element)
	return nil
}

// delayedFocus attempts focus after a short delay
func (fm *FocusManager) delayedFocus() {
	time.Sleep(100 * time.Millisecond)

	fyne.Do(func() {
		if fm.window != nil {
			fm.window.RequestFocus()
		}
	})
}

// delayedFocusElement attempts to focus a specific element after a delay
func (fm *FocusManager) delayedFocusElement(element fyne.Focusable) {
	time.Sleep(100 * time.Millisecond)

	fyne.Do(func() {
		fm.focusDirect(element)
	})
}

// findFocusableElements recursively finds all focusable elements in a container
func (fm *FocusManager) findFocusableElements(obj fyne.CanvasObject) []fyne.Focusable {
	var focusables []fyne.Focusable

	switch v := obj.(type) {
	case fyne.Focusable:
		focusables = append(focusables, v)
	case *fyne.Container:
		for _, child := range v.Objects {
			focusables = append(focusables, fm.findFocusableElements(child)...)
		}
	}
	return focusables
}

// BuildFocusQueue builds the focus queue from the window's content
func (fm *FocusManager) BuildFocusQueue() {
	fyne.Do(func() {
		if fm.window == nil {
			return
		}

		content := fm.window.Content()
		if content == nil {
			return
		}

		focusables := fm.findFocusableElements(content)

		fm.queueMu.Lock()
		fm.focusQueue = focusables
		fm.queueMu.Unlock()

		fm.log.Debug("Focus queue built", "count", len(focusables))
	})
}
