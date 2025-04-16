package Graphics

import "math"

const (
	CLEAR  float32 = 0x0.0p0
	BLACK  float32 = -0x1.0p125
	GRAY   float32 = -0x1.0101p126
	SILVER float32 = -0x1.6d6d6cp126
	WHITE  float32 = -0x1.fffffep126
	RED    float32 = -0x1.0001fep125
	ORANGE float32 = -0x1.00fffep125
	YELLOW float32 = -0x1.01fffep125
	GREEN  float32 = -0x1.01fep125
	BLUE   float32 = -0x1.fep126
)

func ColorFromHEX(hex uint32) float32 {
	alpha := float32((hex>>24)&0xFF) / 255
	red := float32((hex>>16)&0xFF) / 255
	green := float32((hex>>8)&0xFF) / 255
	blue := float32(hex&0xFF) / 255
	return ColorFromRGBA(red, green, blue, alpha)
}

func ColorFromRGBA(r, g, b, a float32) float32 {
	red := uint32(r * 255)
	green := uint32(g * 255)
	blue := uint32(b * 255)
	alpha := uint32(a * 255)
	return math.Float32frombits((alpha << 24) | (blue << 16) | (green << 8) | red)
}

func ColorToHEX(color float32) uint32 {
	r, g, b, a := ColorToRGBA(color)
	return (uint32(a*255) << 24) |
		(uint32(r*255) << 16) |
		(uint32(g*255) << 8) |
		uint32(b*255)
}

func ColorToRGBA(color float32) (r, g, b, a float32) {
	bits := math.Float32bits(color)
	r = float32(bits&0xFF) / 255
	g = float32((bits>>8)&0xFF) / 255
	b = float32((bits>>16)&0xFF) / 255
	a = float32((bits>>24)&0xFF) / 255
	return r, g, b, a
}
