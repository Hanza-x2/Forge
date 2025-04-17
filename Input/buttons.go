package Input

import "github.com/go-gl/glfw/v3.3/glfw"

type MouseButton int

const (
	MouseButton1 MouseButton = iota
	MouseButton2
	MouseButton3
	MouseButton4
	MouseButton5
	MouseButton6
	MouseButton7
	MouseButton8
	MouseButtonLast   = MouseButton8
	MouseButtonLeft   = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMiddle = MouseButton3
)

func convertGLFWButton(button glfw.MouseButton) MouseButton {
	switch button {
	case glfw.MouseButton1:
		return MouseButtonLeft
	case glfw.MouseButton2:
		return MouseButtonRight
	case glfw.MouseButton3:
		return MouseButtonMiddle
	case glfw.MouseButton4:
		return MouseButton4
	case glfw.MouseButton5:
		return MouseButton5
	case glfw.MouseButton6:
		return MouseButton6
	case glfw.MouseButton7:
		return MouseButton7
	case glfw.MouseButton8:
		return MouseButton8
	default:
		return MouseButtonLast
	}
}
