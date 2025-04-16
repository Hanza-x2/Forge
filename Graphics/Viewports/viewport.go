package Viewports

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
)

type Viewport interface {
	Apply(centerCamera bool)

	Update(screenWidth, screenHeight int32, centerCamera bool)

	GetCamera() *Graphics.Camera

	GetWorldDimensions() (float32, float32)

	GetScreenDimensions() (x, y, width, height int32)
}

type BaseViewport struct {
	Camera      *Graphics.Camera
	WorldWidth  float32
	WorldHeight float32

	screenX      int32
	screenY      int32
	screenWidth  int32
	screenHeight int32
}

func (viewport *BaseViewport) GetCamera() *Graphics.Camera {
	return viewport.Camera
}

func (viewport *BaseViewport) GetWorldDimensions() (float32, float32) {
	return viewport.WorldWidth, viewport.WorldHeight
}

func (viewport *BaseViewport) GetScreenDimensions() (x, y, width, height int32) {
	return viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight
}
