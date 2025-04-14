package graphics

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	maxQuads     = 1000
	verticesSize = maxQuads * 4 * (2 + 4 + 2) // x,y + color + u,v
	indicesSize  = maxQuads * 6
)

type Vertex struct {
	Position mgl32.Vec2
	Color    [4]uint8
	TexCoord mgl32.Vec2
}

type Batch struct {
	vao, vbo, ebo uint32
	vertices      []float32
	indices       []uint32
	vertexCount   int
	texture       *Texture
	shader        *Shader
	projection    mgl32.Mat4
	drawing       bool
}

func NewBatch(shader *Shader) *Batch {
	b := &Batch{
		vertices: make([]float32, verticesSize),
		indices:  make([]uint32, indicesSize),
		shader:   shader,
	}

	// Generate indices
	for i, j := 0, 0; i < indicesSize; i, j = i+6, j+4 {
		b.indices[i] = uint32(j)
		b.indices[i+1] = uint32(j + 1)
		b.indices[i+2] = uint32(j + 2)
		b.indices[i+3] = uint32(j + 2)
		b.indices[i+4] = uint32(j + 3)
		b.indices[i+5] = uint32(j)
	}

	// Setup VAO/VBO/EBO
	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.vbo)
	gl.GenBuffers(1, &b.ebo)

	gl.BindVertexArray(b.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(b.vertices)*4, nil, gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(b.indices)*4, gl.Ptr(b.indices), gl.STATIC_DRAW)

	// Position attribute
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	// Color attribute
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.UNSIGNED_BYTE, true, 8*4, gl.PtrOffset(2*4))

	// Texture coordinate attribute
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))

	gl.BindVertexArray(0)

	return b
}

func (b *Batch) Begin() {
	if b.drawing {
		return
	}
	b.drawing = true
	b.vertexCount = 0
}

func (b *Batch) End() {
	if !b.drawing {
		return
	}
	b.Flush()
	b.drawing = false
}

func (b *Batch) Flush() {
	if b.vertexCount == 0 {
		return
	}

	b.shader.Bind()
	b.shader.SetUniformMatrix4fv("u_projection", &b.projection[0])
	b.shader.SetUniform1i("u_texture", 0)

	if b.texture != nil {
		b.texture.Bind(0)
	}

	gl.BindVertexArray(b.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, b.vertexCount*8*4, gl.Ptr(b.vertices))

	gl.DrawElements(gl.TRIANGLES, int32(b.vertexCount/4*6), gl.UNSIGNED_INT, nil)

	b.vertexCount = 0
}

func (b *Batch) SetProjection(projection mgl32.Mat4) {
	b.projection = projection
}

func (b *Batch) Draw(texture *Texture, x, y, width, height float32) {
	if !b.drawing {
		return
	}

	if b.texture != texture {
		b.Flush()
		b.texture = texture
	}

	if b.vertexCount >= maxQuads*4 {
		b.Flush()
	}

	// Define quad vertices (positions, color white, texture coords)
	vertices := []float32{
		x, y, 1, 1, 1, 1, 0, 1,
		x, y + height, 1, 1, 1, 1, 0, 0,
		x + width, y + height, 1, 1, 1, 1, 1, 0,
		x + width, y, 1, 1, 1, 1, 1, 1,
	}

	copy(b.vertices[b.vertexCount*8:], vertices)
	b.vertexCount += 4
}

func (b *Batch) DrawRegion(region *TextureRegion, x, y, width, height float32) {
	if !b.drawing {
		return
	}

	if b.texture != region.Texture {
		b.Flush()
		b.texture = region.Texture
	}

	if b.vertexCount >= maxQuads*4 {
		b.Flush()
	}

	vertices := []float32{
		x, y, 1, 1, 1, 1, region.U, region.V2,
		x, y + height, 1, 1, 1, 1, region.U, region.V,
		x + width, y + height, 1, 1, 1, 1, region.U2, region.V,
		x + width, y, 1, 1, 1, 1, region.U2, region.V2,
	}

	copy(b.vertices[b.vertexCount*8:], vertices)
	b.vertexCount += 4
}

func (b *Batch) Dispose() {
	gl.DeleteVertexArrays(1, &b.vao)
	gl.DeleteBuffers(1, &b.vbo)
	gl.DeleteBuffers(1, &b.ebo)
}
