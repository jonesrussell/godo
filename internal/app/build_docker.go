//go:build docker
// +build docker

package app

// SetDockerEnvironment is called by wire to configure Docker environment
func SetDockerEnvironment(app *App) *App {
	app.SetIsDocker(true)
	return app
}
