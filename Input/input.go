package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type Action int

const (
	ActionRelease Action = iota
	ActionPress
	ActionRepeat
)

type ModifierKey int

const (
	ModShift ModifierKey = 1 << iota
	ModControl
	ModAlt
	ModSuper
)

type KeyCallback func(key Key, action Action, modifiers ModifierKey)
type MouseButtonCallback func(button MouseButton, action Action, modifiers ModifierKey, x, y float32)
type MouseMoveCallback func(x, y float32)
type MouseScrollCallback func(x, y float32)

type Handler struct {
	keys         [KeyLast + 1]bool
	mouseButtons [MouseButtonLast + 1]bool
	mouseX       float32
	mouseY       float32

	keyCallback         KeyCallback
	mouseButtonCallback MouseButtonCallback
	mouseMoveCallback   MouseMoveCallback
	mouseScrollCallback MouseScrollCallback
}

func (handler *Handler) Install(window *glfw.Window) {
	window.SetKeyCallback(handler.keyCallbackWrapper)
	window.SetMouseButtonCallback(handler.mouseButtonCallbackWrapper)
	window.SetCursorPosCallback(handler.mouseMoveCallbackWrapper)
	window.SetScrollCallback(handler.mouseScrollCallbackWrapper)
}

func (handler *Handler) SetKeyCallback(callback KeyCallback) {
	handler.keyCallback = callback
}

func (handler *Handler) SetMouseButtonCallback(callback MouseButtonCallback) {
	handler.mouseButtonCallback = callback
}

func (handler *Handler) SetMouseMoveCallback(callback MouseMoveCallback) {
	handler.mouseMoveCallback = callback
}

func (handler *Handler) SetMouseScrollCallback(callback MouseScrollCallback) {
	handler.mouseScrollCallback = callback
}

func (handler *Handler) IsKeyPressed(key Key) bool {
	if key < 0 || key > KeyLast {
		return false
	}
	return handler.keys[key]
}

func (handler *Handler) IsMouseButtonPressed(button MouseButton) bool {
	if button < 0 || button > MouseButtonLast {
		return false
	}
	return handler.mouseButtons[button]
}

func (handler *Handler) GetMousePosition() (float32, float32) {
	return handler.mouseX, handler.mouseY
}

func (handler *Handler) GetMouseX() float32 {
	return handler.mouseX
}

func (handler *Handler) GetMouseY() float32 {
	return handler.mouseY
}

func convertGLFWAction(glfwAction glfw.Action) Action {
	return Action(glfwAction)
}

func convertGLFWMods(glfwMods glfw.ModifierKey) ModifierKey {
	var mods ModifierKey
	if glfwMods&glfw.ModShift != 0 {
		mods |= ModShift
	}
	if glfwMods&glfw.ModControl != 0 {
		mods |= ModControl
	}
	if glfwMods&glfw.ModAlt != 0 {
		mods |= ModAlt
	}
	if glfwMods&glfw.ModSuper != 0 {
		mods |= ModSuper
	}
	return mods
}

func (handler *Handler) keyCallbackWrapper(_ *glfw.Window, glfwKey glfw.Key, _ int, glfwAction glfw.Action, mods glfw.ModifierKey) {
	key := convertGLFWKey(glfwKey)
	if key >= 0 && key <= KeyLast {
		if glfwAction == glfw.Press {
			handler.keys[key] = true
		} else if glfwAction == glfw.Release {
			handler.keys[key] = false
		}
	}
	if handler.keyCallback != nil {
		handler.keyCallback(key, convertGLFWAction(glfwAction), convertGLFWMods(mods))
	}
}

func (handler *Handler) mouseButtonCallbackWrapper(_ *glfw.Window, glfwButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	button := convertGLFWButton(glfwButton)
	if button >= 0 && int(button) < len(handler.mouseButtons) {
		if action == glfw.Press {
			handler.mouseButtons[button] = true
		} else if action == glfw.Release {
			handler.mouseButtons[button] = false
		}
	}

	if handler.mouseButtonCallback != nil {
		handler.mouseButtonCallback(button, convertGLFWAction(action), convertGLFWMods(mods), handler.mouseX, handler.mouseY)
	}
}

func (handler *Handler) mouseMoveCallbackWrapper(_ *glfw.Window, x, y float64) {
	castX := float32(x)
	castY := float32(y)
	handler.mouseX = castX
	handler.mouseY = castY
	if handler.mouseMoveCallback != nil {
		handler.mouseMoveCallback(castX, castY)
	}
}

func (handler *Handler) mouseScrollCallbackWrapper(_ *glfw.Window, x, y float64) {
	if handler.mouseScrollCallback != nil {
		handler.mouseScrollCallback(float32(x), -float32(y))
	}
}
