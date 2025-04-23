package Viewports

import (
	"github.com/ForgeLeaf/Forge/Graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Viewport interface {
	Apply(centerCamera bool)

	Update(screenWidth, screenHeight float32, centerCamera bool)

	GetCamera() *Graphics.Camera

	GetWorldDimensions() (float32, float32)

	GetScreenDimensions() (x, y, width, height float32)

	Unproject(input mgl32.Vec2) mgl32.Vec2

	Project(input mgl32.Vec2) mgl32.Vec2
}

type BaseViewport struct {
	Camera      *Graphics.Camera
	WorldWidth  float32
	WorldHeight float32

	screenX      float32
	screenY      float32
	screenWidth  float32
	screenHeight float32
}

func (viewport *BaseViewport) Apply(centerCamera bool) {
	gl.Viewport(
		int32(viewport.screenX), int32(viewport.screenY),
		int32(viewport.screenWidth), int32(viewport.screenHeight),
	)
	viewport.Camera.Width = viewport.WorldWidth
	viewport.Camera.Height = viewport.WorldHeight
	if centerCamera {
		viewport.Camera.Position = mgl32.Vec2{viewport.WorldWidth / 2, viewport.WorldHeight / 2}
	}
	viewport.Camera.Update()
}

func (viewport *BaseViewport) GetCamera() *Graphics.Camera {
	return viewport.Camera
}

func (viewport *BaseViewport) GetWorldDimensions() (float32, float32) {
	return viewport.WorldWidth, viewport.WorldHeight
}

func (viewport *BaseViewport) GetScreenDimensions() (x, y, width, height float32) {
	return viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight
}

func (viewport *BaseViewport) Unproject(input mgl32.Vec2) mgl32.Vec2 {
	return viewport.Camera.Unproject(input, viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
}

func (viewport *BaseViewport) Project(input mgl32.Vec2) mgl32.Vec2 {
	return viewport.Camera.Project(input, viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
}
