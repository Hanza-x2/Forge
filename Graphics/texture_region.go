package Graphics

import "math"

type TextureRegion struct {
	Texture *Texture
	U, V    float32
	U2, V2  float32
	Width   int
	Height  int
}

func NewTextureRegion(texture *Texture) *TextureRegion {
	region := &TextureRegion{Texture: texture}
	region.SetBounds(0, 0, texture.Width, texture.Height)
	return region
}

func (r *TextureRegion) SetTexture(texture *Texture) {
	r.Texture = texture
	r.SetBounds(0, 0, texture.Width, texture.Height)
}

func (r *TextureRegion) SetBounds(x, y, width, height int32) {
	invTexWidth := 1.0 / float32(r.Texture.Width)
	invTexHeight := 1.0 / float32(r.Texture.Height)

	r.SetUV(
		float32(x)*invTexWidth,
		float32(y)*invTexHeight,
		float32(x+width)*invTexWidth,
		float32(y+height)*invTexHeight,
	)

	r.Width = int(math.Abs(float64(width)))
	r.Height = int(math.Abs(float64(height)))
}

func (r *TextureRegion) SetUV(u, v, u2, v2 float32) {
	texWidth := float32(r.Texture.Width)
	texHeight := float32(r.Texture.Height)

	r.Width = int(math.Round(math.Abs(float64(u2-u)) * float64(texWidth)))
	r.Height = int(math.Round(math.Abs(float64(v2-v)) * float64(texHeight)))

	if r.Width == 1 && r.Height == 1 {
		adjustX := 0.25 / texWidth
		u += adjustX
		u2 -= adjustX

		adjustY := 0.25 / texHeight
		v += adjustY
		v2 -= adjustY
	}

	r.U = u
	r.V = v
	r.U2 = u2
	r.V2 = v2
}

func (r *TextureRegion) Flip(flipX, flipY bool) {
	if flipX {
		r.U, r.U2 = r.U2, r.U
	}
	if flipY {
		r.V, r.V2 = r.V2, r.V
	}
}
