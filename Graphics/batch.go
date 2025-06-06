package Graphics

import (
	"github.com/ForgeLeaf/Forge"
	"github.com/ForgeLeaf/Forge/Math"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"math"
)

const (
	PositionSize = 2
	ColorSize    = 1
	TexCoordSize = 2
	VertexSize   = PositionSize + ColorSize + TexCoordSize
	Stride       = VertexSize * 4

	maxQuads     = 1000
	verticesSize = maxQuads * 4 * VertexSize
	indicesSize  = maxQuads * 6
)

type Batch struct {
	vao, vbo, ebo uint32
	vertices      []float32
	indices       []uint32
	color         Color
	spaceFactor   float32
	vertexCount   int
	pixel         *TextureRegion
	texture       *Texture
	shader        *Shader
	driver        *Forge.Driver
	projection    mgl32.Mat4
	transform     mgl32.Mat3
	identity      mgl32.Mat3
	drawing       bool
}

func createDefaultShader() *Shader {
	shader, err := NewShader(`
#version 130

in vec4 a_position;
in vec4 a_color;
in vec2 a_texCoord;

uniform mat4 u_projection;

out mediump vec4 v_color;
out highp vec2 v_texCoords;

void main() {
    v_color = a_color;
    v_color.a *= (255.0 / 254.0);
    v_texCoords = a_texCoord;
    gl_Position = u_projection * a_position;
}`, `
#version 130

in mediump vec4 v_color;
in highp vec2 v_texCoords;

uniform highp sampler2D u_texture;

out vec4 fragColor;

void main() {
    fragColor = v_color * texture(u_texture, v_texCoords);
}`)
	if err != nil {
		log.Fatalf("Failed to create default shader: %v", err)
	}
	return shader
}

func createDefaultPixel() *TextureRegion {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	pixel := &Texture{
		ID:     textureID,
		Width:  1,
		Height: 1,
	}

	data := []uint8{255, 255, 255, 255}
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	return NewTextureRegion(pixel)
}

func NewBatch(driver *Forge.Driver) *Batch {
	pixel := createDefaultPixel()
	shader := createDefaultShader()
	batch := &Batch{
		vertices: make([]float32, verticesSize),
		indices:  make([]uint32, indicesSize),
		driver:   driver,
		shader:   shader,
		pixel:    pixel,
	}
	batch.SetColor(ColorWhite)

	for i, j := 0, 0; i < indicesSize; i, j = i+6, j+4 {
		batch.indices[i] = uint32(j)
		batch.indices[i+1] = uint32(j + 1)
		batch.indices[i+2] = uint32(j + 2)
		batch.indices[i+3] = uint32(j + 2)
		batch.indices[i+4] = uint32(j + 3)
		batch.indices[i+5] = uint32(j)
	}

	gl.GenVertexArrays(1, &batch.vao)
	gl.GenBuffers(1, &batch.vbo)
	gl.GenBuffers(1, &batch.ebo)

	gl.BindVertexArray(batch.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(batch.vertices)*4, nil, gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, batch.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(batch.indices)*4, gl.Ptr(batch.indices), gl.STATIC_DRAW)

	posLoc := uint32(gl.GetAttribLocation(shader.Program, gl.Str("a_position\x00")))
	colorLoc := uint32(gl.GetAttribLocation(shader.Program, gl.Str("a_color\x00")))
	texCoordLoc := uint32(gl.GetAttribLocation(shader.Program, gl.Str("a_texCoord\x00")))

	gl.EnableVertexAttribArray(posLoc)
	gl.VertexAttribPointer(posLoc, 2, gl.FLOAT, false, Stride, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(colorLoc)
	gl.VertexAttribPointer(colorLoc, 4, gl.UNSIGNED_BYTE, true, Stride, gl.PtrOffset(2*4)) // After 2 floats

	gl.EnableVertexAttribArray(texCoordLoc)
	gl.VertexAttribPointer(texCoordLoc, 2, gl.FLOAT, false, Stride, gl.PtrOffset(3*4)) // After 2 floats + 4 bytes

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	return batch
}

func (batch *Batch) SetColor(color Color) {
	batch.color = color
}

func (batch *Batch) SetColorRGBA(r, g, b, a float32) {
	batch.SetColor(ColorFromRGBA(r, g, b, a))
}

func (batch *Batch) SetColorHEX(hex uint32) {
	batch.SetColor(ColorFromHEX(hex))
}

func (batch *Batch) GetColor() Color {
	return batch.color
}

func (batch *Batch) Begin() {
	if batch.drawing {
		return
	}
	batch.drawing = true
	batch.vertexCount = 0
}

func (batch *Batch) End() {
	if !batch.drawing {
		return
	}
	batch.Flush()
	batch.drawing = false
}

func (batch *Batch) Flush() {
	if batch.vertexCount == 0 {
		return
	}

	matrix := batch.projection
	if batch.transform != batch.identity {
		matrix = batch.projection.Mul4(Math.Mat3ToMat4(batch.transform))
	}

	batch.shader.Bind()
	batch.shader.SetUniformMatrix4fv("u_projection", &matrix[0])
	batch.shader.SetUniform1i("u_texture", 0)

	if batch.texture != nil {
		batch.texture.Bind(0)
	}

	gl.BindVertexArray(batch.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, batch.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, batch.vertexCount*VertexSize*4, gl.Ptr(batch.vertices))

	gl.DrawElements(gl.TRIANGLES, int32(batch.vertexCount/4*6), gl.UNSIGNED_INT, nil)

	batch.vertexCount = 0
}

func (batch *Batch) SetProjection(projection mgl32.Mat4) {
	if batch.drawing {
		batch.Flush()
	}
	batch.projection = projection
	batch.spaceFactor = 2 / (projection[0] * batch.driver.Width)
}

func (batch *Batch) PushTransform(transform mgl32.Mat3) {
	if batch.drawing {
		batch.Flush()
	}
	batch.transform = transform
}

func (batch *Batch) PopTransform() {
	batch.PushTransform(batch.identity)
}

// Has to be called before doing any drawing *or heavy calculations* (Simple proxies may skip this)
func (batch *Batch) valid() bool {
	if !batch.drawing {
		log.Print("Begin() has to be called before drawing.")
		return false
	}
	return true
}

func (batch *Batch) Push(
	texture *Texture,
	x1, y1 float32, c1 Color, u1, v1 float32,
	x2, y2 float32, c2 Color, u2, v2 float32,
	x3, y3 float32, c3 Color, u3, v3 float32,
	x4, y4 float32, c4 Color, u4, v4 float32,
) {
	if !batch.valid() {
		return
	}

	if batch.texture != texture {
		batch.Flush()
		batch.texture = texture
	}

	if batch.vertexCount >= maxQuads*4 {
		batch.Flush()
	}

	vertices := []float32{
		x1, y1, float32(c1), u1, v1,
		x2, y2, float32(c2), u2, v2,
		x3, y3, float32(c3), u3, v3,
		x4, y4, float32(c4), u4, v4,
	}

	copy(batch.vertices[batch.vertexCount*VertexSize:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) Line(x1, y1, x2, y2 float32, color Color, stroke float32) {
	batch.LineEx(x1, y1, color, x2, y2, color, stroke)
}

func (batch *Batch) LineEx(x1, y1 float32, c1 Color, x2, y2 float32, c2 Color, stroke float32) {
	if !batch.valid() {
		return
	}
	halfStroke := (batch.spaceFactor * stroke) / 2
	dX := x2 - x1
	dY := y2 - y1
	length := float32(math.Sqrt(float64(dX*dX + dY*dY)))
	if length == 0 {
		return
	}
	nX := dY / length * halfStroke
	nY := -dX / length * halfStroke
	batch.FillQuadEx(
		x1-nX, y1-nY, c1,
		x1+nX, y1+nY, c1,
		x2+nX, y2+nY, c2,
		x2-nX, y2-nY, c2,
	)
}

func (batch *Batch) LineRect(x, y, width, height float32, color Color, stroke float32) {
	batch.LineEx(x, y, color, x+width, y, color, stroke)
	batch.LineEx(x+width, y, color, x+width, y+height, color, stroke)
	batch.LineEx(x+width, y+height, color, x, y+height, color, stroke)
	batch.LineEx(x, y+height, color, x, y, color, stroke)
}

func (batch *Batch) LineRectEx(x, y, originX, originY, width, height, scaleX, scaleY, rotation float32, color Color, stroke float32) {
	if !batch.valid() {
		return
	}
	rad := float64(rotation * math.Pi / 180)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))
	fx := -originX
	fy := -originY
	fx2 := width - originX
	fy2 := height - originY
	if scaleX != 1 || scaleY != 1 {
		fx *= scaleX
		fy *= scaleY
		fx2 *= scaleX
		fy2 *= scaleY
	}
	worldOriginX := x + originX
	worldOriginY := y + originY
	x1 := cos*fx - sin*fy + worldOriginX
	y1 := sin*fx + cos*fy + worldOriginY
	x2 := cos*fx2 - sin*fy + worldOriginX
	y2 := sin*fx2 + cos*fy + worldOriginY
	x3 := cos*fx2 - sin*fy2 + worldOriginX
	y3 := sin*fx2 + cos*fy2 + worldOriginY
	x4 := x1 + (x3 - x2)
	y4 := y3 - (y2 - y1)
	batch.LineEx(x1, y1, color, x2, y2, color, stroke)
	batch.LineEx(x2, y2, color, x3, y3, color, stroke)
	batch.LineEx(x3, y3, color, x4, y4, color, stroke)
	batch.LineEx(x4, y4, color, x1, y1, color, stroke)
}

func (batch *Batch) FillQuad(x1, y1, x2, y2, x3, y3, x4, y4 float32, color Color) {
	batch.FillQuadEx(x1, y1, color, x2, y2, color, x3, y3, color, x4, y4, color)
}

func (batch *Batch) FillQuadEx(
	x1, y1 float32, c1 Color,
	x2, y2 float32, c2 Color,
	x3, y3 float32, c3 Color,
	x4, y4 float32, c4 Color,
) {
	batch.Push(batch.pixel.Texture,
		x1, y1, c1, 0, 1,
		x2, y2, c2, 0, 0,
		x3, y3, c3, 1, 0,
		x4, y4, c4, 1, 1,
	)
}

func (batch *Batch) FillRect(x, y, width, height float32, color Color) {
	batch.FillQuad(x, y, x+width, y, x+width, y+height, x, y+height, color)
}

func (batch *Batch) FillRectEx(x, y, width, height float32, c1, c2, c3, c4 Color) {
	batch.FillQuadEx(x, y, c1, x+width, y, c2, x+width, y+height, c3, x, y+height, c4)
}

func (batch *Batch) Draw(texture *Texture, x, y, width, height float32) {
	batch.Push(texture,
		x, y, batch.color, 0, 1,
		x, y+height, batch.color, 0, 0,
		x+width, y+height, batch.color, 1, 0,
		x+width, y, batch.color, 1, 1,
	)
}

func (batch *Batch) DrawEx(
	texture *Texture, x, y, originX, originY,
	width, height, scaleX, scaleY, rotation float32,
	srcX, srcY, srcWidth, srcHeight int,
	flipX, flipY bool,
) {
	if !batch.valid() {
		return
	}
	worldOriginX := x + originX
	worldOriginY := y + originY

	fx := -originX
	fy := -originY
	fx2 := width - originX
	fy2 := height - originY

	if scaleX != 1 || scaleY != 1 {
		fx *= scaleX
		fy *= scaleY
		fx2 *= scaleX
		fy2 *= scaleY
	}

	p1x := fx
	p1y := fy
	p2x := fx
	p2y := fy2
	p3x := fx2
	p3y := fy2
	p4x := fx2
	p4y := fy

	var x1, y1, x2, y2, x3, y3, x4, y4 float32

	if rotation != 0 {
		rad := float64(rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))

		x1 = p1x*cos - p1y*sin
		y1 = p1x*sin + p1y*cos
		x2 = p2x*cos - p2y*sin
		y2 = p2x*sin + p2y*cos
		x3 = p3x*cos - p3y*sin
		y3 = p3x*sin + p3y*cos
		x4 = x1 + (x3 - x2)
		y4 = y3 - (y2 - y1)
	} else {
		x1 = p1x
		y1 = p1y
		x2 = p2x
		y2 = p2y
		x3 = p3x
		y3 = p3y
		x4 = p4x
		y4 = p4y
	}

	x1 += worldOriginX
	x2 += worldOriginX
	x3 += worldOriginX
	x4 += worldOriginX

	y1 += worldOriginY
	y2 += worldOriginY
	y3 += worldOriginY
	y4 += worldOriginY

	invTexWidth := 1.0 / float32(texture.Width)
	invTexHeight := 1.0 / float32(texture.Height)

	u := float32(srcX) * invTexWidth
	v := float32((srcY)+srcHeight) * invTexHeight
	u2 := float32((srcX)+srcWidth) * invTexWidth
	v2 := float32(srcY) * invTexHeight

	if flipX {
		u, u2 = u2, u
	}

	if flipY {
		v, v2 = v2, v
	}

	color := batch.color
	batch.Push(texture,
		x1, y1, color, u, v,
		x2, y2, color, u, v2,
		x3, y3, color, u2, v2,
		x4, y4, color, u2, v,
	)

}

func (batch *Batch) DrawRegion(region *TextureRegion, x, y, width, height float32) {
	batch.Push(region.Texture,
		x, y, batch.color, region.U, region.V2,
		x, y+height, batch.color, region.U, region.V,
		x+width, y+height, batch.color, region.U2, region.V,
		x+width, y, batch.color, region.U2, region.V2,
	)
}

func (batch *Batch) DrawRegionEx(
	region *TextureRegion, x, y, originX, originY,
	width, height, scaleX, scaleY, rotation float32,
) {
	if !batch.valid() {
		return
	}
	worldOriginX := x + originX
	worldOriginY := y + originY

	fx := -originX
	fy := -originY
	fx2 := width - originX
	fy2 := height - originY

	if scaleX != 1 || scaleY != 1 {
		fx *= scaleX
		fy *= scaleY
		fx2 *= scaleX
		fy2 *= scaleY
	}

	p1x := fx
	p1y := fy
	p2x := fx
	p2y := fy2
	p3x := fx2
	p3y := fy2
	p4x := fx2
	p4y := fy

	var x1, y1, x2, y2, x3, y3, x4, y4 float32

	if rotation != 0 {
		rad := float64(rotation * math.Pi / 180)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))

		x1 = p1x*cos - p1y*sin
		y1 = p1x*sin + p1y*cos
		x2 = p2x*cos - p2y*sin
		y2 = p2x*sin + p2y*cos
		x3 = p3x*cos - p3y*sin
		y3 = p3x*sin + p3y*cos
		x4 = x1 + (x3 - x2)
		y4 = y3 - (y2 - y1)
	} else {
		x1 = p1x
		y1 = p1y
		x2 = p2x
		y2 = p2y
		x3 = p3x
		y3 = p3y
		x4 = p4x
		y4 = p4y
	}

	x1 += worldOriginX
	x2 += worldOriginX
	x3 += worldOriginX
	x4 += worldOriginX

	y1 += worldOriginY
	y2 += worldOriginY
	y3 += worldOriginY
	y4 += worldOriginY

	u := region.U
	v := region.V2
	u2 := region.U2
	v2 := region.V

	color := batch.color
	batch.Push(region.Texture,
		x1, y1, color, u, v,
		x2, y2, color, u, v2,
		x3, y3, color, u2, v2,
		x4, y4, color, u2, v,
	)

}

func (batch *Batch) Dispose() {
	gl.DeleteVertexArrays(1, &batch.vao)
	gl.DeleteBuffers(1, &batch.vbo)
	gl.DeleteBuffers(1, &batch.ebo)
	batch.shader.Dispose()
}
