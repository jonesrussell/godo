package gui

import (
	"context"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/logger"
)

func ShowQuickNote(ctx context.Context, application *app.App) {
	logger.Debug("Opening quick note window")

	qnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	quickNote := application.GetQuickNote()
	if quickNote == nil {
		logger.Error("Failed to get quick note instance")
		return
	}

	// Show the quick note window and wait for it to complete
	errChan := make(chan error, 1)
	go func() {
		if err := quickNote.Show(qnCtx); err != nil {
			logger.Error("Failed to show quick note", "error", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Handle input in a separate goroutine
	go func() {
		select {
		case input := <-quickNote.GetInput():
			logger.Debug("Received quick note input", "input", input)
			todoService := application.GetTodoService()
			if _, err := todoService.CreateTodo(qnCtx, input, ""); err != nil {
				logger.Error("Failed to create todo from quick note", "error", err)
			} else {
				logger.Debug("Successfully created todo from quick note")
			}
		case err := <-errChan:
			if err != nil {
				logger.Error("Quick note error", "error", err)
			}
		case <-qnCtx.Done():
			logger.Debug("Quick note context cancelled")
			return
		}
	}()
}
