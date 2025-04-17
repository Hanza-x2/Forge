package Forge

type Configuration interface{}

type WindowConfiguration struct {
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
