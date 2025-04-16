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
	kerning   map[[2]rune]int
	cache     map[rune]*TextureRegion
	cacheLock sync.Mutex
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
		kerning: make(map[[2]rune]int),
		cache:   make(map[rune]*TextureRegion),
	}

	if err := font.parse(file); err != nil {
		return nil, err
	}

	dir := filepath.Dir(path)
	if err := font.loadTextures(dir); err != nil {
		return nil, err
	}

	return font, nil
}

func (f *Font) parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
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
			f.parseInfo(attrs)
		case "common":
			f.parseCommon(attrs)
		case "page":
			if err := f.parsePage(attrs); err != nil {
				return err
			}
		case "char":
			f.parseChar(attrs)
		case "kerning":
			f.parseKerning(attrs)
		}
	}
	return scanner.Err()
}

func (f *Font) parseInfo(attrs map[string]string) {
	f.info.face = attrs["face"]
	f.info.size = atoi(attrs["size"])
	padding := strings.Split(attrs["padding"], ",")
	f.info.padding.up = atoi(padding[0])
	f.info.padding.right = atoi(padding[1])
	f.info.padding.down = atoi(padding[2])
	f.info.padding.left = atoi(padding[3])
}

func (f *Font) parseCommon(attrs map[string]string) {
	f.common.lineHeight = atoi(attrs["lineHeight"])
	f.common.base = atoi(attrs["base"])
	f.common.scaleW = atoi(attrs["scaleW"])
	f.common.scaleH = atoi(attrs["scaleH"])
}

func (f *Font) parsePage(attrs map[string]string) error {
	id := atoi(attrs["id"])
	if id >= len(f.pages) {
		f.pages = append(f.pages, make([]*TextureRegion, id+1-len(f.pages))...)
	}
	return nil
}

func (f *Font) parseChar(attrs map[string]string) {
	id := atoi(attrs["id"])
	f.chars[rune(id)] = struct {
		x, y, width, height int
		xOffset, yOffset    int
		xAdvance            int
		page                int
	}{
		x:        atoi(attrs["x"]),
		y:        atoi(attrs["y"]),
		width:    atoi(attrs["width"]),
		height:   atoi(attrs["height"]),
		xOffset:  atoi(attrs["xoffset"]),
		yOffset:  atoi(attrs["yoffset"]),
		xAdvance: atoi(attrs["xadvance"]),
		page:     atoi(attrs["page"]),
	}
}

func (f *Font) parseKerning(attrs map[string]string) {
	first := rune(atoi(attrs["first"]))
	second := rune(atoi(attrs["second"]))
	f.kerning[[2]rune{first, second}] = atoi(attrs["amount"])
}

func (f *Font) loadTextures(textureDir string) error {
	for i := range f.pages {
		texture, err := NewTexture(filepath.Join(textureDir, f.getPageFile(i)))
		if err != nil {
			return err
		}
		f.pages[i] = NewTextureRegion(texture)
	}
	return nil
}

func (f *Font) getPageFile(id int) string {
	return strconv.Itoa(id) + ".png"
}

func (f *Font) getGlyph(ch rune) *TextureRegion {
	f.cacheLock.Lock()
	defer f.cacheLock.Unlock()

	if region, exists := f.cache[ch]; exists {
		return region
	}

	char, exists := f.chars[ch]
	if !exists || char.page >= len(f.pages) || f.pages[char.page] == nil {
		return nil
	}

	page := f.pages[char.page]
	region := &TextureRegion{
		Texture: page.Texture,
		U:       float32(char.x) / float32(f.common.scaleW),
		V:       float32(char.y) / float32(f.common.scaleH),
		U2:      float32(char.x+char.width) / float32(f.common.scaleW),
		V2:      float32(char.y+char.height) / float32(f.common.scaleH),
		Width:   char.width,
		Height:  char.height,
	}

	f.cache[ch] = region
	return region
}

func (f *Font) Draw(batch *Batch, text string, x, y float32) {
	f.DrawEx(batch, text, x, y, 0, len(text), -1, 0, false, "")
}

func (f *Font) DrawEx(batch *Batch, text string, x, y float32, start, end int, targetWidth float32, hAlign int, wrap bool, truncate string) {
	if !batch.valid() || start >= end {
		return
	}

	if end > len(text) {
		end = len(text)
	}

	textWidth := f.MeasureText(text[start:end])
	offsetX := float32(0)
	if targetWidth > 0 {
		switch hAlign {
		case 1: // AlignCenter
			offsetX = (targetWidth - textWidth) / 2
		case 2: // AlignRight
			offsetX = targetWidth - textWidth
		}
	}

	currentX, currentY := x+offsetX, y
	for i := start; i < end; i++ {
		ch := rune(text[i])
		if ch == '\n' {
			currentX = x + offsetX
			currentY += float32(f.common.lineHeight)
			continue
		}

		char, exists := f.chars[ch]
		if !exists {
			continue
		}

		kerning := 0
		if i > start {
			kerning = f.kerning[[2]rune{rune(text[i-1]), ch}]
		}

		if region := f.getGlyph(ch); region != nil {
			posX := currentX + float32(char.xOffset+kerning)
			posY := currentY + float32(f.common.base-char.yOffset-char.height)
			batch.DrawRegion(region, posX, posY, float32(region.Width), float32(region.Height))
		}

		currentX += float32(char.xAdvance + kerning)
	}
}

func (f *Font) MeasureText(text string) float32 {
	width, maxWidth := float32(0), float32(0)
	for i, ch := range text {
		if ch == '\n' {
			if width > maxWidth {
				maxWidth = width
			}
			width = 0
			continue
		}

		char, exists := f.chars[ch]
		if !exists {
			continue
		}

		kerning := 0
		if i > 0 {
			kerning = f.kerning[[2]rune{rune(text[i-1]), ch}]
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
