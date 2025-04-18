package Geometry

type Rectangle struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func (rectangle *Rectangle) Set(x, y, width, height float32) {
	rectangle.X = x
	rectangle.Y = y
	rectangle.Width = width
	rectangle.Height = height
}

func (rectangle *Rectangle) SetPosition(x, y float32) {
	rectangle.X = x
	rectangle.Y = y
}

func (rectangle *Rectangle) SetSize(width, height float32) {
	rectangle.Width = width
	rectangle.Height = height
}

func (rectangle *Rectangle) Contains(x, y float32) bool {
	return x >= rectangle.X && x <= rectangle.X+rectangle.Width &&
		y >= rectangle.Y && y <= rectangle.Y+rectangle.Height
}

func (rectangle *Rectangle) Overlaps(other *Rectangle) bool {
	return rectangle.X < other.X+other.Width && rectangle.X+rectangle.Width > other.X &&
		rectangle.Y < other.Y+other.Height && rectangle.Y+rectangle.Height > other.Y
}
