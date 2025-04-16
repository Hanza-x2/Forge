package Geometry

type Shape interface {
	Contains(x, y float32) bool
}
