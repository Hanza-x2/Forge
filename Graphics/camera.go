package Graphics

import (
	"forgejo.max7.fun/m.alkhatib/GoForge/Math"
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

func (camera *Camera) Unproject(input mgl32.Vec2, viewportX, viewportY, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	return Math.Vec2MulMat4(mgl32.Vec2{
		(2*(input.X()-viewportX))/viewportWidth - 1,
		(2*(input.Y()-viewportY))/viewportHeight - 1,
	}, camera.Inverse)
}

func (camera *Camera) Project(input mgl32.Vec2, viewportX, viewportY, viewportWidth, viewportHeight float32) mgl32.Vec2 {
	output := Math.Vec2MulMat4(input, camera.Matrix)
	return mgl32.Vec2{
		viewportWidth*(output.X()+1)/2 + viewportX,
		viewportHeight*(output.Y()+1)/2 + viewportY,
	}
}
