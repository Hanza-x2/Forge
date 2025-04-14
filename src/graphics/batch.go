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
	batch.projection = projection
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

	// Define quad vertices (positions, color white, texture coords)
	vertices := []float32{
		x, y, 1, 1, 1, 1, 0, 1,
		x, y + height, 1, 1, 1, 1, 0, 0,
		x + width, y + height, 1, 1, 1, 1, 1, 0,
		x + width, y, 1, 1, 1, 1, 1, 1,
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
		x, y, 1, 1, 1, 1, region.U, region.V2,
		x, y + height, 1, 1, 1, 1, region.U, region.V,
		x + width, y + height, 1, 1, 1, 1, region.U2, region.V,
		x + width, y, 1, 1, 1, 1, region.U2, region.V2,
	}

	copy(batch.vertices[batch.vertexCount*8:], vertices)
	batch.vertexCount += 4
}

func (batch *Batch) Dispose() {
	gl.DeleteVertexArrays(1, &batch.vao)
	gl.DeleteBuffers(1, &batch.vbo)
	gl.DeleteBuffers(1, &batch.ebo)
}
