package Viewports

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Graphics"
)

type FitViewport struct {
	BaseViewport
}

func NewFitViewport(worldWidth, worldHeight float32, screenWidth, screenHeight float32) *FitViewport {
	viewport := &FitViewport{
		BaseViewport: BaseViewport{
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			Camera:      Graphics.NewCamera(worldWidth, worldHeight),
		},
	}
	viewport.Update(screenWidth, screenHeight, true)
	return viewport
}

func (viewport *FitViewport) Update(screenWidth, screenHeight float32, centerCamera bool) {
	targetRatio := screenHeight / screenWidth
	sourceRatio := viewport.WorldHeight / viewport.WorldWidth

	scale := targetRatio
	if targetRatio > sourceRatio {
		scale = screenWidth / viewport.WorldWidth
	} else {
		scale = screenHeight / viewport.WorldHeight
	}

	scaledWidth := viewport.WorldWidth * scale
	scaledHeight := viewport.WorldHeight * scale

	viewport.screenX = (screenWidth - scaledWidth) / 2
	viewport.screenY = (screenHeight - scaledHeight) / 2
	viewport.screenWidth = scaledWidth
	viewport.screenHeight = scaledHeight

	viewport.Apply(centerCamera)
}
