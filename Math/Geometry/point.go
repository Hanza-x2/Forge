package Geometry

type Point struct {
	X float32
	Y float32
}

func (point *Point) Set(x, y float32) {
	point.X = x
	point.Y = y
}

func (point *Point) Contains(x, y float32) bool {
	return point.X == x && point.Y == y
}
