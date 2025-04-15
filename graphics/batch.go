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
	batch.SetColor(1, 1, 1, 1)

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

func (batch *Batch) SetColor(r, g, b, a float32) {
	red := uint32(r * 255)
	green := uint32(g * 255)
	blue := uint32(b * 255)
	alpha := uint32(a * 255)
	batch.color = math.Float32frombits((alpha << 24) | (blue << 16) | (green << 8) | red)
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
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, batch.vertexCount*8*4, gl.Ptr(batch.vertices))

	gl.DrawElements(gl.TRIANGLES, int32(batch.vertexCount/4*6), gl.UNSIGNED_INT, nil)

	batch.vertexCount = 0
}

func (batch *Batch) SetProjection(projection mgl32.Mat4) {
	if !batch.drawing {
		batch.Flush()
	}
	batch.projection = projection
}

func (batch *Batch) FillQuad(x, y, x2, y2 float32) {
	color := batch.color
	batch.FillQuadEx(x, y, x2, y2, color, color, color, color)
}

func (batch *Batch) FillQuadEx(x, y, x2, y2, c1, c2, c3, c4 float32) {
	if !batch.drawing {
		return
	}

	if batch.texture != batch.pixel.Texture {
		batch.Flush()
		batch.texture = batch.pixel.Texture
	}

	if batch.vertexCount >= maxQuads*4 {
		batch.Flush()
	}

	vertices := []float32{
		x, y, c1, 0, 1,
		x, y2, c2, 0, 0,
		x2, y2, c3, 1, 0,
		x2, y, c4, 1, 1,
	}

	copy(batch.vertices[batch.vertexCount*8:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) FillRect(x, y, width, height float32) {
	color := batch.color
	batch.FillRectEx(x, y, width, height, color, color, color, color)
}

func (batch *Batch) FillRectEx(x, y, width, height, c1, c2, c3, c4 float32) {
	batch.FillQuadEx(x, y, x+width, y+height, c1, c2, c3, c4)
}

func (batch *Batch) Draw(texture *Texture, x, y, width, height float32) {
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
		x, y, batch.color, 0, 1,
		x, y + height, batch.color, 0, 0,
		x + width, y + height, batch.color, 1, 0,
		x + width, y, batch.color, 1, 1,
	}

	copy(batch.vertices[batch.vertexCount*8:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) DrawRegion(region *TextureRegion, x, y, width, height float32) {
	if !batch.drawing {
		return
	}

	if batch.texture != region.Texture {
		batch.Flush()
		batch.texture = region.Texture
	}

	if batch.vertexCount >= maxQuads*4 {
		batch.Flush()
	}

	vertices := []float32{
		x, y, batch.color, region.U, region.V2,
		x, y + height, batch.color, region.U, region.V,
		x + width, y + height, batch.color, region.U2, region.V,
		x + width, y, batch.color, region.U2, region.V2,
	}

	copy(batch.vertices[batch.vertexCount*8:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) Dispose() {
	gl.DeleteVertexArrays(1, &batch.vao)
	gl.DeleteBuffers(1, &batch.vbo)
	gl.DeleteBuffers(1, &batch.ebo)
	batch.shader.Dispose()
}
