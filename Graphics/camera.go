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

	camera.Matrix = mgl32.Ortho(
		camera.Position.X()-width/2,
		camera.Position.X()+width/2,
		camera.Position.Y()-height/2,
		camera.Position.Y()+height/2,
		-1, 1,
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

func (camera *Camera) Unproject(input mgl32.Vec2, viewportX, viewportY, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	x := (2*(input.X()-viewportX))/viewportWidth - 1
	y := 1 - (2*(input.Y()-viewportY))/viewportHeight
	output := mgl32.TransformCoordinate(mgl32.Vec3{x, y, 0}, camera.Inverse)
	return mgl32.Vec2{output.X(), output.Y()}
}

func (camera *Camera) Project(input mgl32.Vec2, viewportX, viewportY, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	output := mgl32.TransformCoordinate(mgl32.Vec3{input.X(), input.Y(), 0}, camera.Matrix)
	return mgl32.Vec2{
		viewportWidth*(output.X()+1)/2 + viewportX,
		viewportHeight*(1-output.Y())/2 + viewportY,
	}
}
