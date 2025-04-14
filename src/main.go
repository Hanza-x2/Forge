package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

type App struct {
	driver *Driver
}

func (app *App) Create(driver *Driver) {
	app.driver = driver
	log.Println("Application created")
}

func (app *App) Render(driver *Driver, delta float32) {

	if driver.Input.IsKeyPressed(glfw.KeyEscape) {
		log.Println("Escape key pressed, exiting...")
		driver.Stop()
	}

	gl.ClearColor(0.1, 0.2, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (app *App) Resize(driver *Driver, width, height int) {
	log.Printf("Application resized to %d x %d", width, height)
}

func (app *App) Destroy(driver *Driver) {
	log.Println("Application destroyed")
}

func init() {
	runtime.LockOSThread()
}

func main() {
	configuration := DefaultConfig()
	if err := RunSafe(&App{}, configuration); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
