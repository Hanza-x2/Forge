package Graphics

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Font struct {
	info struct {
		face    string
		size    int
		padding struct{ up, right, down, left int }
	}
	common struct {
		lineHeight, base, scaleW, scaleH int
	}
	pages []*TextureRegion
	chars map[rune]struct {
		x, y, width, height int
		xOffset, yOffset    int
		xAdvance            int
		page                int
	}
	kerning    map[[2]rune]int
	cache      map[rune]*TextureRegion
	cacheLock  sync.Mutex
	textureDir string
}

func LoadFont(path string) (*Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	font := &Font{
		chars: make(map[rune]struct {
			x, y, width, height int
			xOffset, yOffset    int
			xAdvance            int
			page                int
		}),
		kerning:    make(map[[2]rune]int),
		cache:      make(map[rune]*TextureRegion),
		textureDir: filepath.Dir(path),
	}

	if err := font.parse(file); err != nil {
		return nil, err
	}

	return font, nil
}

func (font *Font) parse(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}

		attrs := make(map[string]string)
		for _, part := range parts[1:] {
			if kv := strings.Split(part, "="); len(kv) == 2 {
				attrs[kv[0]] = strings.Trim(kv[1], `"`)
			}
		}

		switch parts[0] {
		case "info":
			font.parseInfo(attrs)
		case "common":
			font.parseCommon(attrs)
		case "page":
			if err := font.parsePage(attrs); err != nil {
				return err
			}
		case "char":
			font.parseChar(attrs)
		case "kerning":
			font.parseKerning(attrs)
		}
	}
	return scanner.Err()
}

func (font *Font) parseInfo(attributes map[string]string) {
	font.info.face = attributes["face"]
	font.info.size = atoi(attributes["size"])
	padding := strings.Split(attributes["padding"], ",")
	font.info.padding.up = atoi(padding[0])
	font.info.padding.right = atoi(padding[1])
	font.info.padding.down = atoi(padding[2])
	font.info.padding.left = atoi(padding[3])
}

func (font *Font) parseCommon(attributes map[string]string) {
	font.common.lineHeight = atoi(attributes["lineHeight"])
	font.common.base = atoi(attributes["base"])
	font.common.scaleW = atoi(attributes["scaleW"])
	font.common.scaleH = atoi(attributes["scaleH"])
}

func (font *Font) parsePage(attributes map[string]string) error {
	id := atoi(attributes["id"])
	if id >= len(font.pages) {
		font.pages = append(font.pages, make([]*TextureRegion, id+1-len(font.pages))...)
	}
	texture, err := NewTexture(filepath.Join(font.textureDir, attributes["file"]))
	if err != nil {
		return err
	}
	font.pages[id] = NewTextureRegion(texture)
	return nil
}

func (font *Font) parseChar(attributes map[string]string) {
	id := atoi(attributes["id"])
	font.chars[rune(id)] = struct {
		x, y, width, height int
		xOffset, yOffset    int
		xAdvance            int
		page                int
	}{
		x:        atoi(attributes["x"]),
		y:        atoi(attributes["y"]),
		width:    atoi(attributes["width"]),
		height:   atoi(attributes["height"]),
		xOffset:  atoi(attributes["xoffset"]),
		yOffset:  atoi(attributes["yoffset"]),
		xAdvance: atoi(attributes["xadvance"]),
		page:     atoi(attributes["page"]),
	}
}

func (font *Font) parseKerning(attributes map[string]string) {
	first := rune(atoi(attributes["first"]))
	second := rune(atoi(attributes["second"]))
	font.kerning[[2]rune{first, second}] = atoi(attributes["amount"])
}

func (font *Font) getGlyph(ch rune) *TextureRegion {
	font.cacheLock.Lock()
	defer font.cacheLock.Unlock()

	if region, exists := font.cache[ch]; exists {
		return region
	}

	char, exists := font.chars[ch]
	if !exists || char.page >= len(font.pages) || font.pages[char.page] == nil {
		return nil
	}

	page := font.pages[char.page]
	region := &TextureRegion{
		Texture: page.Texture,
		U:       float32(char.x) / float32(font.common.scaleW),
		V:       float32(char.y) / float32(font.common.scaleH),
		U2:      float32(char.x+char.width) / float32(font.common.scaleW),
		V2:      float32(char.y+char.height) / float32(font.common.scaleH),
		Width:   char.width,
		Height:  char.height,
	}

	font.cache[ch] = region
	return region
}

func (font *Font) Draw(batch *Batch, text string, x, y float32) {
	font.DrawEx(batch, text, x, y, 0, len(text), -1, 0, false, "")
}

func (font *Font) DrawEx(batch *Batch, text string, x, y float32, start, end int, targetWidth float32, hAlign int, wrap bool, truncate string) {
	if !batch.valid() || start >= end {
		return
	}

	if end > len(text) {
		end = len(text)
	}

	textWidth := font.MeasureText(text[start:end])
	var offsetX float32 = 0
	if targetWidth > 0 && hAlign > 0 {
		if hAlign == 1 {
			offsetX = (targetWidth - textWidth) / 2
		} else if hAlign == 2 {
			offsetX = targetWidth - textWidth
		}
	}

	currentX, currentY := x+offsetX, y
	for i := start; i < end; i++ {
		ch := rune(text[i])
		if ch == '\n' {
			currentX = x + offsetX
			currentY += float32(font.common.lineHeight)
			continue
		}

		char, exists := font.chars[ch]
		if !exists {
			continue
		}

		kerning := 0
		if i > start {
			kerning = font.kerning[[2]rune{rune(text[i-1]), ch}]
		}

		if region := font.getGlyph(ch); region != nil {
			posX := currentX + float32(char.xOffset+kerning)
			posY := currentY + float32(font.common.base-char.yOffset-char.height)
			batch.DrawRegion(region, posX, posY, float32(region.Width), float32(region.Height))
		}

		currentX += float32(char.xAdvance + kerning)
	}
}

func (font *Font) MeasureText(text string) float32 {
	width, maxWidth := float32(0), float32(0)
	for i, ch := range text {
		if ch == '\n' {
			if width > maxWidth {
				maxWidth = width
			}
			width = 0
			continue
		}

		char, exists := font.chars[ch]
		if !exists {
			continue
		}

		kerning := 0
		if i > 0 {
			kerning = font.kerning[[2]rune{rune(text[i-1]), ch}]
		}

		if i == len(text)-1 {
			width += float32(char.xOffset + char.width + kerning)
		} else {
			width += float32(char.xAdvance + kerning)
		}
	}

	if width > maxWidth {
		return width
	}
	return maxWidth
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
