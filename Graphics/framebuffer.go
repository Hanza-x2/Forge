package Graphics

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type FrameBuffer struct {
	width, height int
	fbo           uint32
	colorTexture  Texture
}

func NewFrameBuffer(width, height int) (*FrameBuffer, error) {
	frameBuffer := &FrameBuffer{
		width:  width,
		height: height,
	}

	gl.GenFramebuffers(1, &frameBuffer.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer.fbo)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureID, 0)
	frameBuffer.colorTexture = Texture{
		ID:     textureID,
		Width:  width,
		Height: height,
	}

	if status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); status != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("framebuffer incomplete: 0x%x", status)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	return frameBuffer, nil
}

func (frameBuffer *FrameBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer.fbo)
	gl.Viewport(0, 0, int32(frameBuffer.width), int32(frameBuffer.height))
}

func (frameBuffer *FrameBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (frameBuffer *FrameBuffer) GetColorTexture() Texture {
	return frameBuffer.colorTexture
}

func (frameBuffer *FrameBuffer) Dispose() {
	frameBuffer.colorTexture.Dispose()
	gl.DeleteFramebuffers(1, &frameBuffer.fbo)
}
