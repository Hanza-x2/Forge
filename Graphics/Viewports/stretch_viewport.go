package Viewports

import (
	"github.com/ForgeLeaf/Forge/Graphics"
)

type StretchViewport struct {
	BaseViewport
}

func NewStretchViewport(worldWidth, worldHeight float32, screenWidth, screenHeight float32) *StretchViewport {
	viewport := &StretchViewport{
		BaseViewport: BaseViewport{
			WorldWidth:   worldWidth,
			WorldHeight:  worldHeight,
			Camera:       Graphics.NewCamera(worldWidth, worldHeight),
			screenX:      0,
			screenY:      0,
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
	}
	viewport.Update(screenWidth, screenHeight, true)
	return viewport
}

func (viewport *StretchViewport) Update(screenWidth, screenHeight float32, centerCamera bool) {
	viewport.screenX = 0
	viewport.screenY = 0
	viewport.screenWidth = screenWidth
	viewport.screenHeight = screenHeight
	viewport.Apply(centerCamera)
}
