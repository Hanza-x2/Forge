package Graphics

import (
	"errors"
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/scanner"
)

type Font struct {
	Descriptor *Descriptor
	Pages      map[int]*TextureRegion
	glyphCache map[rune]*TextureRegion
	cacheMutex sync.Mutex
}

type Descriptor struct {
	Info    Info
	Common  Common
	Pages   map[int]Page
	Chars   map[rune]Char
	Kerning map[CharPair]Kerning
}

type Info struct {
	Face     string
	Size     int
	Bold     bool
	Italic   bool
	Charset  string
	Unicode  bool
	StretchH int
	Smooth   bool
	AA       int
	Padding  Padding
	Spacing  Spacing
	Outline  int
}

type Padding struct {
	Up, Right, Down, Left int
}

type Spacing struct {
	Horizontal, Vertical int
}

type Common struct {
	LineHeight   int
	Base         int
	ScaleW       int
	ScaleH       int
	Packed       bool
	AlphaChannel int
	RedChannel   int
	GreenChannel int
	BlueChannel  int
}

type Page struct {
	ID   int
	File string
}

type CharPair struct {
	First, Second rune
}

type Char struct {
	ID       rune
	X        int
	Y        int
	Width    int
	Height   int
	XOffset  int
	YOffset  int
	XAdvance int
	Page     int
	Channel  int
}

type Kerning struct {
	Amount int
}

func (char *Char) Pos() image.Point {
	return image.Pt(char.X, char.Y)
}

func (char *Char) Size() image.Point {
	return image.Pt(char.Width, char.Height)
}

func (char *Char) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: char.Pos(),
		Max: char.Pos().Add(char.Size()),
	}
}

func (char *Char) Offset() image.Point {
	return image.Pt(char.XOffset, char.YOffset)
}

func LoadDescriptor(path string) (descriptor *Descriptor, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeChecked(f, &err)
	return parseDescriptor(filepath.Base(path), f)
}

func ReadDescriptor(reader io.Reader) (descriptor *Descriptor, err error) {
	return parseDescriptor("bmfont", reader)
}

func parseDescriptor(filename string, reader io.Reader) (*Descriptor, error) {
	var parser tagsParser
	tags, err := parser.parse(filename, reader)
	if err != nil {
		return nil, err
	}
	font := Descriptor{
		Pages:   make(map[int]Page),
		Chars:   make(map[rune]Char),
		Kerning: make(map[CharPair]Kerning),
	}
	for _, tag := range tags {
		switch tag.name {
		case "info":
			var values = tag.intListAttr("padding", 4)
			var values2 = tag.intListAttr("spacing", 2)
			font.Info = Info{
				Face:     tag.stringAttr("face"),
				Size:     tag.intAttr("size"),
				Bold:     tag.boolAttr("bold"),
				Italic:   tag.boolAttr("italic"),
				Charset:  tag.stringAttr("charset"),
				Unicode:  tag.boolAttr("unicode"),
				StretchH: tag.intAttr("stretchH"),
				Smooth:   tag.boolAttr("smooth"),
				AA:       tag.intAttr("aa"),
				Padding: Padding{
					Up:    values[0],
					Right: values[1],
					Down:  values[2],
					Left:  values[3],
				},
				Spacing: Spacing{
					Horizontal: values2[0],
					Vertical:   values2[1],
				},
				Outline: tag.intAttr("outline"),
			}
		case "common":
			font.Common = Common{
				LineHeight:   tag.intAttr("lineHeight"),
				Base:         tag.intAttr("base"),
				ScaleW:       tag.intAttr("scaleW"),
				ScaleH:       tag.intAttr("scaleH"),
				Packed:       tag.boolAttr("packed"),
				AlphaChannel: tag.intAttr("alphaChnl"),
				RedChannel:   tag.intAttr("redChnl"),
				GreenChannel: tag.intAttr("greenChnl"),
				BlueChannel:  tag.intAttr("blueChnl"),
			}
		case "page":
			id := tag.intAttr("id")
			font.Pages[id] = Page{
				ID:   id,
				File: tag.stringAttr("file"),
			}
		case "char":
			id := tag.runeAttr("id")
			font.Chars[id] = Char{
				ID:       id,
				X:        tag.intAttr("x"),
				Y:        tag.intAttr("y"),
				Width:    tag.intAttr("width"),
				Height:   tag.intAttr("height"),
				XOffset:  tag.intAttr("xoffset"),
				YOffset:  tag.intAttr("yoffset"),
				XAdvance: tag.intAttr("xadvance"),
			}
		case "kerning":
			pair := CharPair{
				First:  tag.runeAttr("first"),
				Second: tag.runeAttr("second"),
			}
			font.Kerning[pair] = Kerning{
				Amount: tag.intAttr("amount"),
			}
		}
	}
	return &font, nil
}

type tagsParser struct {
	errors  errorList
	scanner scanner.Scanner
	pos     scanner.Position
	tok     rune
	lit     string
}

func (parser *tagsParser) next() {
	parser.tok = parser.scanner.Scan()
	parser.pos = parser.scanner.Position
	parser.lit = parser.scanner.TokenText()
}

func (parser *tagsParser) parse(filename string, reader io.Reader) ([]tag, error) {
	parser.scanner.Init(reader)
	parser.scanner.Filename = filename
	parser.scanner.Whitespace ^= 1 << '\n'
	parser.scanner.Error = func(s *scanner.Scanner, msg string) {}
	parser.next()

	var tags []tag
	for parser.tok != scanner.EOF {
		tagName := parser.lit
		parser.expect(scanner.Ident, "tag name")
		attrs := make(map[string]string)
		for parser.tok != '\n' && parser.tok != scanner.EOF {
			attrName := parser.lit
			parser.expect(scanner.Ident, "attribute name")
			value := ""
			var err error
			parser.expect('=', `"="`)
			switch parser.tok {
			case scanner.String:
				value, err = strconv.Unquote(parser.lit)
				if err != nil {
					parser.tok = '\n'
					continue
				}
				if parser.scanner.Peek() == '"' {
					parser.scanner.Next()
					value += `"`
				}
				parser.next()
			case scanner.Int, '-':
				value = parser.parseIntList()
			default:
				parser.errorExpected("string or integer attribute value")
			}
			attrs[attrName] = value
		}
		tags = append(tags, tag{
			name:  tagName,
			attrs: attrs,
		})
		parser.next()
	}
	return tags, parser.errors.Err()
}

func (parser *tagsParser) parseIntList() string {
	var sb strings.Builder
	for parser.tok == scanner.Int || parser.tok == ',' || parser.tok == '-' {
		sb.WriteString(parser.lit)
		parser.next()
	}
	return sb.String()
}

func (parser *tagsParser) expect(tok rune, msg string) {
	if parser.tok != tok {
		parser.errorExpected(msg)
	}
	parser.next()
}

func (parser *tagsParser) errorExpected(msg string) {
	parser.error(newError(parser.pos, "expected "+msg+", found "+scanner.TokenString(parser.tok)))
}

func (parser *tagsParser) error(err error) {
	parser.errors = append(parser.errors, err)
}

func newError(pos scanner.Position, msg string) error {
	return errors.New(pos.String() + ": " + msg)
}

type tag struct {
	name  string
	attrs map[string]string
}

func (t *tag) intAttr(name string) int {
	value, _ := strconv.Atoi(t.stringAttr(name))
	return value
}

func (t *tag) runeAttr(name string) rune {
	value, _ := strconv.ParseInt(t.stringAttr(name), 10, 32)
	return rune(value)
}

func (t *tag) stringAttr(name string) string {
	return t.attrs[name]
}

func (t *tag) boolAttr(name string) bool {
	return t.intAttr(name) != 0
}

func (t *tag) intListAttr(name string, n int) []int {
	values := make([]int, n)
	parts := strings.Split(t.stringAttr(name), ",")
	for i, part := range parts {
		if i == len(values) {
			break
		}
		value, _ := strconv.Atoi(strings.TrimSpace(part))
		values[i] = value
	}
	return values
}

type errorList []error

func (list errorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

func (list errorList) Error() string {
	if len(list) == 0 {
		return "no errors"
	}
	return list[0].Error()
}

func LoadFont(path string) (font *Font, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeChecked(file, &err)
	dir, _ := filepath.Split(path)
	return ReadFont(file, dir)
}

func ReadFont(reader io.Reader, textureDir string) (*Font, error) {
	desc, err := ReadDescriptor(reader)
	if err != nil {
		return nil, err
	}
	font := Font{
		Descriptor: desc,
		Pages:      make(map[int]*TextureRegion),
		glyphCache: make(map[rune]*TextureRegion),
	}

	for id, page := range desc.Pages {
		texture, err := NewTexture(filepath.Join(textureDir, page.File))
		if err != nil {
			return nil, err
		}

		font.Pages[id] = NewTextureRegion(texture)
	}

	return &font, nil
}

func (font *Font) getGlyphRegion(ch rune) *TextureRegion {
	font.cacheMutex.Lock()
	defer font.cacheMutex.Unlock()

	if region, exists := font.glyphCache[ch]; exists {
		return region
	}

	char, exists := font.Descriptor.Chars[ch]
	if !exists {
		return nil
	}

	page, exists := font.Pages[char.Page]
	if !exists {
		return nil
	}

	region := &TextureRegion{
		Texture: page.Texture,
		U:       float32(char.X) / float32(font.Descriptor.Common.ScaleW),
		V:       float32(char.Y) / float32(font.Descriptor.Common.ScaleH),
		U2:      float32(char.X+char.Width) / float32(font.Descriptor.Common.ScaleW),
		V2:      float32(char.Y+char.Height) / float32(font.Descriptor.Common.ScaleH),
		Width:   char.Width,
		Height:  char.Height,
	}

	font.glyphCache[ch] = region
	return region
}

func closeChecked(c io.Closer, err *error) {
	cErr := c.Close()
	if cErr != nil && *err == nil {
		*err = cErr
	}
}

const (
	AlignLeft = iota
	AlignCenter
	AlignRight
)

func (font *Font) Draw(batch *Batch, text string, x, y float32) {
	font.DrawEx(batch, text, x, y, 0, len(text), -1, AlignLeft, false, "")
}

func (font *Font) DrawEx(batch *Batch, text string, x, y float32, start, end int, targetWidth float32, hAlign int, wrap bool, truncate string) {
	if !batch.valid() {
		return
	}

	if end > len(text) {
		end = len(text)
	}
	if start >= end {
		return
	}

	textWidth := font.MeasureText(text[start:end])
	offsetX := float32(0)
	if targetWidth > 0 {
		switch hAlign {
		case AlignCenter:
			offsetX = (targetWidth - textWidth) / 2
		case AlignRight:
			offsetX = targetWidth - textWidth
		}
	}

	currentX := x + offsetX
	currentY := y

	for i := start; i < end; i++ {
		ch := rune(text[i])

		if ch == '\n' {
			currentX = x + offsetX
			currentY += float32(font.Descriptor.Common.LineHeight)
			continue
		}

		char, exists := font.Descriptor.Chars[ch]
		if !exists {
			continue
		}

		var kerning int
		if i > start {
			prevCh := rune(text[i-1])
			if k, ok := font.Descriptor.Kerning[CharPair{prevCh, ch}]; ok {
				kerning = k.Amount
			}
		}

		charRegion := font.getGlyphRegion(ch)
		if charRegion == nil {
			continue
		}

		posX := currentX + float32(char.XOffset+kerning)
		posY := currentY + float32(font.Descriptor.Common.Base-char.YOffset-char.Height)
		batch.DrawRegion(charRegion, posX, posY, float32(char.Width), float32(char.Height))

		currentX += float32(char.XAdvance + kerning)
	}
}

func (font *Font) MeasureText(text string) float32 {
	width := float32(0)
	maxWidth := float32(0)

	for i, ch := range text {
		if ch == '\n' {
			if width > maxWidth {
				maxWidth = width
			}
			width = 0
			continue
		}

		char, exists := font.Descriptor.Chars[ch]
		if !exists {
			continue
		}

		var kerning int
		if i > 0 {
			prevCh := rune(text[i-1])
			if k, ok := font.Descriptor.Kerning[CharPair{prevCh, ch}]; ok {
				kerning = k.Amount
			}
		}

		if i == len(text)-1 {
			width += float32(char.XOffset + char.Width + kerning)
		} else {
			width += float32(char.XAdvance + kerning)
		}
	}

	if width > maxWidth {
		maxWidth = width
	}
	return maxWidth
}
