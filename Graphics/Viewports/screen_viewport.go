package Viewports

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ScreenViewport struct {
	BaseViewport
}

func NewScreenViewport(screenWidth, screenHeight int) *ScreenViewport {
	worldWidth := float32(screenWidth)
	worldHeight := float32(screenHeight)
	viewport := &ScreenViewport{
		BaseViewport: BaseViewport{
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			Camera:      Graphics.NewCamera(worldWidth, worldHeight),
		},
	}
	viewport.Update(screenWidth, screenHeight, true)
	return viewport
}

func (viewport *ScreenViewport) Apply(centerCamera bool) {
	gl.Viewport(viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
	viewport.Camera.Width = viewport.WorldWidth
	viewport.Camera.Height = viewport.WorldHeight
	if centerCamera {
		viewport.Camera.Position = mgl32.Vec2{viewport.WorldWidth / 2, viewport.WorldHeight / 2}
	}
	viewport.Camera.Update()
}

func (viewport *ScreenViewport) Update(screenWidth, screenHeight int, centerCamera bool) {
	viewport.WorldWidth = float32(screenWidth)
	viewport.WorldHeight = float32(screenHeight)
	viewport.Apply(centerCamera)
}
