package Viewports

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ScreenViewport struct {
	BaseViewport
}

func NewScreenViewport(screenWidth, screenHeight int32) *ScreenViewport {
	viewport := &ScreenViewport{
		BaseViewport: BaseViewport{
			WorldWidth:   float32(screenWidth),
			WorldHeight:  float32(screenHeight),
			Camera:       Graphics.NewCamera(float32(screenWidth), float32(screenHeight)),
			screenX:      0,
			screenY:      0,
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
	}
	return viewport
}

func (viewport *ScreenViewport) Apply(centerCamera bool) {
	gl.Viewport(viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
	viewport.Camera.Width = float32(viewport.screenWidth)
	viewport.Camera.Height = float32(viewport.screenHeight)
	if centerCamera {
		viewport.Camera.Position = mgl32.Vec2{
			float32(viewport.screenWidth) / 2,
			float32(viewport.screenHeight) / 2,
		}
	}
	viewport.Camera.Update()
}

func (viewport *ScreenViewport) Update(screenWidth, screenHeight int32, centerCamera bool) {
	viewport.screenX = 0
	viewport.screenY = 0
	viewport.screenWidth = screenWidth
	viewport.screenHeight = screenHeight
	viewport.WorldWidth = float32(screenWidth)
	viewport.WorldHeight = float32(screenHeight)
	viewport.Apply(centerCamera)
}
