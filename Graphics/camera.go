package Graphics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position mgl32.Vec2
	Zoom     float32
	Width    float32
	Height   float32
	Matrix   mgl32.Mat4
	Inverse  mgl32.Mat4
}

func NewCamera(width, height float32) *Camera {
	cam := &Camera{
		Width:    width,
		Height:   height,
		Zoom:     1.0,
		Position: mgl32.Vec2{width / 2, height / 2},
	}
	cam.Update()
	return cam
}

func (camera *Camera) Update() {
	width := camera.Width * camera.Zoom
	height := camera.Height * camera.Zoom

	camera.Matrix = mgl32.Ortho2D(
		camera.Position.X()-width/2,
		camera.Position.X()+width/2,
		camera.Position.Y()-height/2,
		camera.Position.Y()+height/2,
	)
	camera.Inverse = camera.Matrix.Inv()
}

func (camera *Camera) Resize(width, height float32) {
	camera.Width = width
	camera.Height = height
	camera.Update()
}

func (camera *Camera) Translate(x, y float32) {
	camera.Position = camera.Position.Add(mgl32.Vec2{x, y})
	camera.Update()
}

func (camera *Camera) Unproject(screenPos mgl32.Vec2, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	normalizedX := 2*screenPos.X()/viewportWidth - 1
	normalizedY := 1 - 2*screenPos.Y()/viewportHeight

	worldPos := camera.Inverse.Mul4x1(mgl32.Vec4{normalizedX, normalizedY, 0, 1})
	return mgl32.Vec2{worldPos.X(), worldPos.Y()}.Add(camera.Position)
}

func (camera *Camera) Project(worldPos mgl32.Vec2, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	relativePos := worldPos.Sub(camera.Position)
	clipPos := camera.Matrix.Mul4x1(mgl32.Vec4{relativePos.X(), relativePos.Y(), 0, 1})

	return mgl32.Vec2{
		(clipPos.X() + 1) * viewportWidth / 2,
		(1 - clipPos.Y()) * viewportHeight / 2,
	}
}
