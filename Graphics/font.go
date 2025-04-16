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

type textLine struct {
	text  string
	width float32
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

	lines := font.breakText(text[start:end], targetWidth, wrap, truncate)
	if len(lines) == 0 {
		return
	}

	lineHeight := float32(font.common.lineHeight)
	startY := y - float32(font.common.base)

	for i, line := range lines {
		lineWidth := font.MeasureText(line.text)
		offsetX := float32(0)
		if targetWidth > 0 {
			switch hAlign {
			case 1: // AlignCenter
				offsetX = (targetWidth - lineWidth) / 2
			case 2: // AlignRight
				offsetX = targetWidth - lineWidth
			}
		}

		currentX := x + offsetX
		currentY := startY + float32(i)*lineHeight

		for j, ch := range line.text {
			char, exists := font.chars[ch]
			if !exists {
				continue
			}

			kerning := 0
			if j > 0 {
				kerning = font.kerning[[2]rune{rune(line.text[j-1]), ch}]
			}

			if region := font.getGlyph(ch); region != nil {
				posX := currentX + float32(char.xOffset+kerning)
				posY := currentY + float32(char.yOffset+char.height)
				batch.DrawRegion(region, posX, posY, float32(region.Width), float32(region.Height))
			}

			currentX += float32(char.xAdvance + kerning)
		}
	}
}

func (font *Font) breakText(text string, targetWidth float32, wrap bool, truncate string) []textLine {
	if targetWidth <= 0 || (!wrap && truncate == "") {
		return []textLine{{text: text, width: font.MeasureText(text)}}
	}

	var lines []textLine
	words := strings.Fields(text)

	if truncate != "" {
		fullWidth := font.MeasureText(text)
		if fullWidth <= targetWidth {
			return []textLine{{text: text, width: fullWidth}}
		}

		truncWidth := font.MeasureText(truncate)
		availableWidth := targetWidth - truncWidth

		lastFit := 0
		currentWidth := float32(0)
		for i, ch := range text {
			char, exists := font.chars[ch]
			if !exists {
				continue
			}

			kerning := 0
			if i > 0 {
				kerning = font.kerning[[2]rune{rune(text[i-1]), ch}]
			}

			charWidth := float32(char.xAdvance + kerning)
			if currentWidth+charWidth > availableWidth {
				break
			}

			currentWidth += charWidth
			lastFit = i + 1
		}

		if lastFit > 0 {
			truncated := text[:lastFit] + truncate
			return []textLine{{text: truncated, width: currentWidth + truncWidth}}
		}
		return []textLine{{text: truncate, width: truncWidth}}
	}

	currentLine := ""
	currentWidth := float32(0)
	spaceWidth := float32(font.chars[' '].xAdvance)

	for _, word := range words {
		wordWidth := font.MeasureText(word)
		if currentWidth+wordWidth <= targetWidth {
			if currentLine != "" {
				currentLine += " "
				currentWidth += spaceWidth
			}
			currentLine += word
			currentWidth += wordWidth
		} else {
			if currentLine != "" {
				lines = append(lines, textLine{text: currentLine, width: currentWidth})
			}
			currentLine = word
			currentWidth = wordWidth
		}
	}

	if currentLine != "" {
		lines = append(lines, textLine{text: currentLine, width: currentWidth})
	}

	return lines
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
