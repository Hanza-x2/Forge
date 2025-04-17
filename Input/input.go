package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type KeyCallback func(key glfw.Key, action glfw.Action)
type MouseButtonCallback func(button glfw.MouseButton, action glfw.Action)
type MouseMoveCallback func(x, y float32)
type MouseScrollCallback func(x, y float32)

type Handler struct {
	keys         [glfw.KeyLast + 1]bool
	mouseButtons [glfw.MouseButtonLast + 1]bool
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

func (handler *Handler) IsKeyPressed(key glfw.Key) bool {
	if key < 0 || key > glfw.KeyLast {
		return false
	}
	return handler.keys[key]
}

func (handler *Handler) IsMouseButtonPressed(button glfw.MouseButton) bool {
	if button < 0 || button > glfw.MouseButtonLast {
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

func (handler *Handler) keyCallbackWrapper(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
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

func (handler *Handler) mouseButtonCallbackWrapper(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
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
