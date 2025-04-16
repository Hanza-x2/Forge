package Geometry

type Point struct {
	X float32
	Y float32
}

func (point *Point) Contains(x, y float32) bool {
	return point.X == x && point.Y == y
}
