package viewports

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type StretchViewport struct {
	BaseViewport
}

func NewStretchViewport(worldWidth, worldHeight float32, screenWidth, screenHeight int) *StretchViewport {
	viewport := &StretchViewport{
		BaseViewport: BaseViewport{
			WorldWidth:   worldWidth,
			WorldHeight:  worldHeight,
			Camera:       graphics.NewCamera(worldWidth, worldHeight),
			screenX:      0,
			screenY:      0,
			screenWidth:  int32(screenWidth),
			screenHeight: int32(screenHeight),
		},
	}
	viewport.Update(screenWidth, screenHeight, true)
	return viewport
}

func (viewport *StretchViewport) Apply(centerCamera bool) {
	gl.Viewport(viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
	viewport.Camera.Width = viewport.WorldWidth
	viewport.Camera.Height = viewport.WorldHeight
	if centerCamera {
		viewport.Camera.Position = mgl32.Vec2{viewport.WorldWidth / 2, viewport.WorldHeight / 2}
	}
	viewport.Camera.Update()
}

func (viewport *StretchViewport) Update(screenWidth, screenHeight int, centerCamera bool) {
	viewport.screenX = 0
	viewport.screenY = 0
	viewport.screenWidth = int32(screenWidth)
	viewport.screenHeight = int32(screenHeight)
	viewport.Apply(centerCamera)
}
