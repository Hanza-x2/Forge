package Geometry

type Circle struct {
	X      float32
	Y      float32
	Radius float32
}

func (circle *Circle) Set(x, y, radius float32) {
	circle.X = x
	circle.Y = y
	circle.Radius = radius
}

func (circle *Circle) SetPosition(x, y float32) {
	circle.X = x
	circle.Y = y
}

func (circle *Circle) SetRadius(radius float32) {
	circle.Radius = radius
}

func (circle *Circle) Contains(x, y float32) bool {
	return (x-circle.X)*(x-circle.X)+(y-circle.Y)*(y-circle.Y) <= circle.Radius*circle.Radius
}

func (circle *Circle) Overlaps(other *Circle) bool {
	return (circle.X-other.X)*(circle.X-other.X)+(circle.Y-other.Y)*(circle.Y-other.Y) <= (circle.Radius+other.Radius)*(circle.Radius+other.Radius)
}
