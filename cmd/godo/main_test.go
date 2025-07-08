//go:build !docker && wireinject

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/jonesrussell/godo/internal/app"
)

func TestMainFlow(t *testing.T) {
	t.Run("successful flow", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := app.NewMockApplication(ctrl)

		// Set up expectations
		mockApp.EXPECT().SetupUI().Return(nil)
		mockApp.EXPECT().Run()

		// Execute the flow
		err := mockApp.SetupUI()
		require.NoError(t, err)

		mockApp.Run()
	})

	t.Run("setup failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := app.NewMockApplication(ctrl)

		// Set up expectations for setup failure
		mockApp.EXPECT().SetupUI().Return(errors.New("setup failed"))

		// Execute the flow
		err := mockApp.SetupUI()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "setup failed")
	})

	t.Run("run after successful setup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := app.NewMockApplication(ctrl)

		// Set up expectations
		mockApp.EXPECT().SetupUI().Return(nil)
		mockApp.EXPECT().Run()

		// Execute the flow
		err := mockApp.SetupUI()
		require.NoError(t, err)

		mockApp.Run()
	})

	t.Run("cleanup", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := app.NewMockApplication(ctrl)

		// Set up expectations
		mockApp.EXPECT().Cleanup()

		// Execute cleanup
		mockApp.Cleanup()
	})

	t.Run("full application lifecycle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockApp := app.NewMockApplication(ctrl)

		// Set up expectations for full lifecycle
		mockApp.EXPECT().SetupUI().Return(nil)
		mockApp.EXPECT().Run()
		mockApp.EXPECT().Cleanup()

		// Execute full lifecycle
		err := mockApp.SetupUI()
		require.NoError(t, err)

		mockApp.Run()
		mockApp.Cleanup()
	})
}
