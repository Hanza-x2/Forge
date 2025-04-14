package graphics

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	ID     uint32
	Width  int
	Height int
}

func NewTexture(filePath string) (*Texture, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found: %v", filePath, err)
	}
	defer func(imgFile *os.File) {
		err := imgFile.Close()
		if err != nil {
			panic("failed to close image file")
		}
	}(imgFile)

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode texture: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, errors.New("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&rgba.Pix[0]),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return &Texture{
		ID:     textureID,
		Width:  rgba.Rect.Size().X,
		Height: rgba.Rect.Size().Y,
	}, nil
}

func (texture *Texture) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
}

func (texture *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (texture *Texture) Dispose() {
	gl.DeleteTextures(1, &texture.ID)
}
