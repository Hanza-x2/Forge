package graphics

import (
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
	color         float32
	vertexCount   int
	pixel         *TextureRegion
	texture       *Texture
	shader        *Shader
	projection    mgl32.Mat4
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

func NewBatch() *Batch {
	pixel := createDefaultPixel()
	shader := createDefaultShader()
	batch := &Batch{
		vertices: make([]float32, verticesSize),
		indices:  make([]uint32, indicesSize),
		shader:   shader,
		pixel:    pixel,
	}
	batch.SetColor(WHITE)

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

func (batch *Batch) SetColor(color float32) {
	batch.color = color
}

func (batch *Batch) SetColorRGBA(r, g, b, a float32) {
	batch.SetColor(ColorFromRGBA(r, g, b, a))
}

func (batch *Batch) SetColorHEX(hex uint32) {
	batch.SetColor(ColorFromHEX(hex))
}

func (batch *Batch) GetColor() float32 {
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

	batch.shader.Bind()
	batch.shader.SetUniformMatrix4fv("u_projection", &batch.projection[0])
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
	if !batch.drawing {
		batch.Flush()
	}
	batch.projection = projection
}

func (batch *Batch) Push(
	texture *Texture,
	x1, y1, c1, u1, v1,
	x2, y2, c2, u2, v2,
	x3, y3, c3, u3, v3,
	x4, y4, c4, u4, v4 float32,
) {
	if !batch.drawing {
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
		x1, y1, c1, u1, v1,
		x2, y2, c2, u2, v2,
		x3, y3, c3, u3, v3,
		x4, y4, c4, u4, v4,
	}

	copy(batch.vertices[batch.vertexCount*VertexSize:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) FillQuad(x1, y1, x2, y2, x3, y3, x4, y4 float32) {
	color := batch.color
	batch.FillQuadEx(x1, y1, color, x2, y2, color, x3, y3, color, x4, y4, color)
}

func (batch *Batch) FillQuadEx(x1, y1, c1, x2, y2, c2, x3, y3, c3, x4, y4, c4 float32) {
	batch.Push(batch.pixel.Texture,
		x1, y1, c1, 0, 1,
		x2, y2, c2, 0, 0,
		x3, y3, c3, 1, 0,
		x4, y4, c4, 1, 1,
	)
}

func (batch *Batch) FillRect(x, y, width, height float32) {
	batch.FillQuad(x, y, x+width, y, x+width, y+height, x, y+height)
}

func (batch *Batch) FillRectEx(x, y, width, height, c1, c2, c3, c4 float32) {
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
		cos := float32(math.Cos(float64(rotation)))
		sin := float32(math.Sin(float64(rotation)))

		x1 = p1x*cos - p1y*sin
		y1 = p1x*sin + p1y*cos
		x2 = p2x*cos - p2y*sin
		y2 = p3x*sin + p2y*cos
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

func (batch *Batch) Dispose() {
	gl.DeleteVertexArrays(1, &batch.vao)
	gl.DeleteBuffers(1, &batch.vbo)
	gl.DeleteBuffers(1, &batch.ebo)
	batch.shader.Dispose()
}
