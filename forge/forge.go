package forge

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"time"
)

type WindowConfiguration struct {
	Title                   string
	Width                   int
	Height                  int
	Resizable               bool
	Decorated               bool
	OpenGLVersionMajor      int
	OpenGLVersionMinor      int
	OpenGLProfile           int
	OpenGLForwardCompatible bool
	TargetFPS               int
}

type Application interface {
	Create(driver *Driver)
	Render(driver *Driver, delta float32)
	Resize(driver *Driver, width, height int)
	Destroy(driver *Driver)
}

type Driver struct {
	Input               *InputHandler
	Width               int
	Height              int
	App                 Application
	configuration       WindowConfiguration
	running             bool
	lastFrame           time.Time
	fps                 int
	frames              int
	startTime           time.Time
	targetFrameDuration time.Duration
}
type KeyCallback func(key glfw.Key, action glfw.Action)
type MouseButtonCallback func(button glfw.MouseButton, action glfw.Action)
type MouseMoveCallback func(x, y float64)
type MouseScrollCallback func(x, y float64)

type InputHandler struct {
	keys         [glfw.KeyLast + 1]bool
	mouseButtons [glfw.MouseButtonLast + 1]bool
	mouseX       float64
	mouseY       float64

	keyCallback         KeyCallback
	mouseButtonCallback MouseButtonCallback
	mouseMoveCallback   MouseMoveCallback
	mouseScrollCallback MouseScrollCallback
}

func (handler *InputHandler) Install(window *glfw.Window) {
	window.SetKeyCallback(handler.keyCallbackWrapper)
	window.SetMouseButtonCallback(handler.mouseButtonCallbackWrapper)
	window.SetCursorPosCallback(handler.mouseMoveCallbackWrapper)
	window.SetScrollCallback(handler.mouseScrollCallbackWrapper)
}

func (handler *InputHandler) SetKeyCallback(callback KeyCallback) {
	handler.keyCallback = callback
}

func (handler *InputHandler) SetMouseButtonCallback(callback MouseButtonCallback) {
	handler.mouseButtonCallback = callback
}

func (handler *InputHandler) SetMouseMoveCallback(callback MouseMoveCallback) {
	handler.mouseMoveCallback = callback
}

func (handler *InputHandler) SetMouseScrollCallback(callback MouseScrollCallback) {
	handler.mouseScrollCallback = callback
}

func (handler *InputHandler) IsKeyPressed(key glfw.Key) bool {
	if key < 0 || key > glfw.KeyLast {
		return false
	}
	return handler.keys[key]
}

func (handler *InputHandler) IsMouseButtonPressed(button glfw.MouseButton) bool {
	if button < 0 || button > glfw.MouseButtonLast {
		return false
	}
	return handler.mouseButtons[button]
}

func (handler *InputHandler) GetMousePosition() (float64, float64) {
	return handler.mouseX, handler.mouseY
}

func (handler *InputHandler) GetMouseX() float64 {
	return handler.mouseX
}

func (handler *InputHandler) GetMouseY() float64 {
	return handler.mouseY
}

func (handler *InputHandler) keyCallbackWrapper(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	if key >= 0 && int(key) < len(handler.keys) {
		if action == glfw.Press {
			handler.keys[key] = true
		} else if action == glfw.Release {
			handler.keys[key] = false
		}
	}

	if handler.keyCallback != nil {
		handler.keyCallback(key, action)
	}
}

func (handler *InputHandler) mouseButtonCallbackWrapper(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	if button >= 0 && int(button) < len(handler.mouseButtons) {
		if action == glfw.Press {
			handler.mouseButtons[button] = true
		} else if action == glfw.Release {
			handler.mouseButtons[button] = false
		}
	}

	if handler.mouseButtonCallback != nil {
		handler.mouseButtonCallback(button, action)
	}
}

func (handler *InputHandler) mouseMoveCallbackWrapper(_ *glfw.Window, x, y float64) {
	handler.mouseX, handler.mouseY = x, y

	if handler.mouseMoveCallback != nil {
		handler.mouseMoveCallback(x, y)
	}
}

func (handler *InputHandler) mouseScrollCallbackWrapper(_ *glfw.Window, x, y float64) {
	if handler.mouseScrollCallback != nil {
		handler.mouseScrollCallback(x, y)
	}
}

func glfwBool(value bool) int {
	if value {
		return glfw.True
	}
	return glfw.False
}

func DefaultConfig() WindowConfiguration {
	return WindowConfiguration{
		Title:                   "Leaf",
		Width:                   800,
		Height:                  600,
		Resizable:               true,
		Decorated:               true,
		OpenGLVersionMajor:      3,
		OpenGLVersionMinor:      3,
		OpenGLProfile:           glfw.OpenGLCoreProfile,
		OpenGLForwardCompatible: true,
		TargetFPS:               60,
	}
}

func RunSafe(application Application, configuration WindowConfiguration) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	driver := CreateDriver(application, configuration)
	driver.Start()
	return nil
}

func CreateDriver(application Application, configuration WindowConfiguration) *Driver {
	return &Driver{
		App:                 application,
		Input:               &InputHandler{},
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

	window, err := glfw.CreateWindow(configuration.Width, configuration.Height, configuration.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		driver.Width = width
		driver.Height = height
		driver.App.Resize(driver, width, height)
	})

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	width, height := window.GetSize()
	driver.Width = width
	driver.Height = height

	driver.App.Create(driver)
	driver.App.Resize(driver, width, height)

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

func (driver *Driver) GetFPS() int {
	return driver.fps
}

func (driver *Driver) GetFrames() int {
	return driver.frames
}

func (driver *Driver) IsRunning() bool {
	return driver.running
}
