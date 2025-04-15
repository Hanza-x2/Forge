package viewports

import (
	"GoForge/graphics"
)

type Viewport interface {
	Apply(centerCamera bool)

	Update(screenWidth, screenHeight int, centerCamera bool)

	GetCamera() *graphics.Camera

	GetWorldDimensions() (float32, float32)

	GetScreenDimensions() (x, y, width, height int32)
}

type BaseViewport struct {
	Camera      *graphics.Camera
	WorldWidth  float32
	WorldHeight float32

	screenX      int32
	screenY      int32
	screenWidth  int32
	screenHeight int32
}

func (viewport *BaseViewport) GetCamera() *graphics.Camera {
	return viewport.Camera
}

func (viewport *BaseViewport) GetWorldDimensions() (float32, float32) {
	return viewport.WorldWidth, viewport.WorldHeight
}

func (viewport *BaseViewport) GetScreenDimensions() (x, y, width, height int32) {
	return viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight
}
