package Viewports

import (
	"github.com/ForgeLeaf/Forge/Graphics"
)

type ScreenViewport struct {
	BaseViewport
}

func NewScreenViewport(screenWidth, screenHeight float32) *ScreenViewport {
	viewport := &ScreenViewport{
		BaseViewport: BaseViewport{
			WorldWidth:   screenWidth,
			WorldHeight:  screenHeight,
			Camera:       Graphics.NewCamera(screenWidth, screenHeight),
			screenX:      0,
			screenY:      0,
			screenWidth:  screenWidth,
			screenHeight: screenHeight,
		},
	}
	return viewport
}

func (viewport *ScreenViewport) Update(screenWidth, screenHeight float32, centerCamera bool) {
	viewport.screenX = 0
	viewport.screenY = 0
	viewport.screenWidth = screenWidth
	viewport.screenHeight = screenHeight
	viewport.WorldWidth = screenWidth
	viewport.WorldHeight = screenHeight
	viewport.Apply(centerCamera)
}
