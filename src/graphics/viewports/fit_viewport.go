package viewports

import (
	"GoForge/graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type FitViewport struct {
	BaseViewport
}

func NewFitViewport(worldWidth, worldHeight float32, screenWidth, screenHeight int) *FitViewport {
	viewport := &FitViewport{
		BaseViewport: BaseViewport{
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			Camera:      graphics.NewCamera(worldWidth, worldHeight),
		},
	}
	viewport.Update(screenWidth, screenHeight, true)
	return viewport
}

func (viewport *FitViewport) Apply(centerCamera bool) {
	gl.Viewport(viewport.screenX, viewport.screenY, viewport.screenWidth, viewport.screenHeight)
	viewport.Camera.Width = viewport.WorldWidth
	viewport.Camera.Height = viewport.WorldHeight
	if centerCamera {
		viewport.Camera.Position = mgl32.Vec2{viewport.WorldWidth / 2, viewport.WorldHeight / 2}
	}
	viewport.Camera.Update()
}

func (viewport *FitViewport) Update(screenWidth, screenHeight int, centerCamera bool) {
	targetRatio := float32(screenHeight) / float32(screenWidth)
	sourceRatio := viewport.WorldHeight / viewport.WorldWidth

	scale := targetRatio
	if targetRatio > sourceRatio {
		scale = float32(screenWidth) / viewport.WorldWidth
	} else {
		scale = float32(screenHeight) / viewport.WorldHeight
	}

	scaledWidth := int32(viewport.WorldWidth * scale)
	scaledHeight := int32(viewport.WorldHeight * scale)

	viewport.screenX = (int32(screenWidth) - scaledWidth) / 2
	viewport.screenY = (int32(screenHeight) - scaledHeight) / 2
	viewport.screenWidth = scaledWidth
	viewport.screenHeight = scaledHeight

	viewport.Apply(centerCamera)
}
