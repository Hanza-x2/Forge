package Forge

import "github.com/go-gl/glfw/v3.3/glfw"

type Configuration interface{}

type DesktopConfiguration struct {
	Title                   string
	Width                   int32
	Height                  int32
	Resizable               bool
	Decorated               bool
	OpenGLVersionMajor      int
	OpenGLVersionMinor      int
	OpenGLProfile           int
	OpenGLForwardCompatible bool
	TargetFPS               int32
}

func DefaultDesktopConfig() DesktopConfiguration {
	return DesktopConfiguration{
		Title:                   "Forge Application",
		Width:                   800,
		Height:                  450,
		Resizable:               true,
		Decorated:               true,
		OpenGLVersionMajor:      3,
		OpenGLVersionMinor:      3,
		OpenGLProfile:           glfw.OpenGLCoreProfile,
		OpenGLForwardCompatible: true,
		TargetFPS:               60,
	}
}
