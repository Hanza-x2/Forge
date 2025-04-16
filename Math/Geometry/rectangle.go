package Geometry

type Rectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func (rectangle *Rectangle) Contains(x, y float32) bool {
	return x >= rectangle.X && x <= rectangle.X+rectangle.Width &&
		y >= rectangle.Y && y <= rectangle.Y+rectangle.Height
}

func (rectangle *Rectangle) Overlaps(other *Rectangle) bool {
	return rectangle.X < other.X+other.Width && rectangle.X+rectangle.Width > other.X &&
		rectangle.Y < other.Y+other.Height && rectangle.Y+rectangle.Height > other.Y
}
