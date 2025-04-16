package Geometry

type Ellipse struct {
	X       float32
	Y       float32
	RadiusX float32
	RadiusY float32
}

func (ellipse *Ellipse) Contains(x, y float32) bool {
	x -= ellipse.X
	y -= ellipse.Y
	return (x*x)/(ellipse.RadiusX/2*ellipse.RadiusX/2)+
		(y*y)/(ellipse.RadiusY/2*ellipse.RadiusY/2) <= 1
}
