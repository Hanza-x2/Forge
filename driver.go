package Forge

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Input"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"time"
)

type Driver struct {
	Input               *Input.Handler
	Width               float32
	Height              float32
	App                 Application
	configuration       DesktopConfiguration
	running             bool
	lastFrame           time.Time
	fps                 int32
	frames              int32
	startTime           time.Time
	targetFrameDuration time.Duration
}

func glfwBool(value bool) int {
	if value {
		return glfw.True
	}
	return glfw.False
}

func RunSafe(application Application, configuration DesktopConfiguration) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	driver := CreateDriver(application, configuration)
	driver.Start()
	return nil
}

func CreateDriver(application Application, configuration DesktopConfiguration) *Driver {
	return &Driver{
		App:                 application,
		Input:               &Input.Handler{},
		configuration:       configuration,
		targetFrameDuration: time.Second / time.Duration(configuration.TargetFPS),
		startTime:           time.Now(),
		lastFrame:           time.Now(),
	}
}

func (driver *Driver) Start() {
	if driver.App == nil {
		panic("Application not initialized")
	}

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	configuration := driver.configuration
	glfw.WindowHint(glfw.Resizable, glfwBool(configuration.Resizable))
	glfw.WindowHint(glfw.Decorated, glfwBool(configuration.Decorated))
	glfw.WindowHint(glfw.ContextVersionMajor, configuration.OpenGLVersionMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, configuration.OpenGLVersionMinor)
	glfw.WindowHint(glfw.OpenGLProfile, configuration.OpenGLProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfwBool(configuration.OpenGLForwardCompatible))

	window, err := glfw.CreateWindow(int(configuration.Width), int(configuration.Height), configuration.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		castWidth := float32(width)
		castHeight := float32(height)
		driver.Width = castWidth
		driver.Height = castHeight
		driver.App.Resize(driver, castWidth, castHeight)
	})

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	width, height := window.GetSize()
	castWidth := float32(width)
	castHeight := float32(height)
	driver.Width = castWidth
	driver.Height = castHeight

	driver.App.Create(driver)
	driver.App.Resize(driver, castWidth, castHeight)

	driver.running = true
	driver.Input.Install(window)

	for !window.ShouldClose() && driver.running {
		currentTime := time.Now()
		elapsed := currentTime.Sub(driver.lastFrame)

		if elapsed >= driver.targetFrameDuration {
			delta := float32(elapsed.Seconds())
			driver.lastFrame = currentTime
			driver.frames++

			if time.Since(driver.startTime) >= time.Second {
				driver.fps = driver.frames
				driver.frames = 0
				driver.startTime = currentTime
			}

			glfw.PollEvents()
			driver.App.Render(driver, delta)
			window.SwapBuffers()
		} else {
			time.Sleep(driver.targetFrameDuration - elapsed)
		}
	}

	driver.App.Destroy(driver)
	glfw.Terminate()
}

func (driver *Driver) Stop() {
	driver.running = false
}

func (driver *Driver) GetFPS() int32 {
	return driver.fps
}

func (driver *Driver) GetFrames() int32 {
	return driver.frames
}

func (driver *Driver) IsRunning() bool {
	return driver.running
}
