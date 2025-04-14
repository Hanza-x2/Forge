package viewports

import (
	"GoLeaf/graphics"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type FitViewport struct {
	Camera      *graphics.Camera
	WorldWidth  float32
	WorldHeight float32

	screenX      int32
	screenY      int32
	screenWidth  int32
	screenHeight int32
}

func NewFitViewport(worldWidth, worldHeight float32, screenWidth, screenHeight int) *FitViewport {
	vp := &FitViewport{
		WorldWidth:  worldWidth,
		WorldHeight: worldHeight,
		Camera:      graphics.NewCamera(worldWidth, worldHeight),
	}
	vp.Update(screenWidth, screenHeight, true)
	return vp
}

func (v *FitViewport) Apply(centerCamera bool) {
	gl.Viewport(v.screenX, v.screenY, v.screenWidth, v.screenHeight)
	v.Camera.Width = v.WorldWidth
	v.Camera.Height = v.WorldHeight
	if centerCamera {
		v.Camera.Position = mgl32.Vec2{v.WorldWidth / 2, v.WorldHeight / 2}
	}
	v.Camera.Update()
}

func (v *FitViewport) Update(screenWidth, screenHeight int, centerCamera bool) {
	targetRatio := float32(screenHeight) / float32(screenWidth)
	sourceRatio := v.WorldHeight / v.WorldWidth

	scale := targetRatio
	if targetRatio > sourceRatio {
		scale = float32(screenWidth) / v.WorldWidth
	} else {
		scale = float32(screenHeight) / v.WorldHeight
	}

	scaledWidth := int32(v.WorldWidth * scale)
	scaledHeight := int32(v.WorldHeight * scale)

	v.screenX = (int32(screenWidth) - scaledWidth) / 2
	v.screenY = (int32(screenHeight) - scaledHeight) / 2
	v.screenWidth = scaledWidth
	v.screenHeight = scaledHeight

	v.Apply(centerCamera)
}
