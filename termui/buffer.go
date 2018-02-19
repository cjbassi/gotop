package termui

import "image"

// Cell is a rune with assigned Fg and Bg
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}

// Buffer is a renderable rectangle cell data container.
type Buffer struct {
	Area    image.Rectangle // selected drawing area
	CellMap map[image.Point]Cell
}

func NewCell(ch rune, Fg, Bg Attribute) Cell {
	return Cell{ch, Fg, Bg}
}

// NewBuffer returns a new Buffer
func NewBuffer() *Buffer {
	return &Buffer{
		CellMap: make(map[image.Point]Cell),
		Area:    image.Rectangle{}}
}

// NewFilledBuffer returns a new Buffer filled with ch, fb and bg.
func NewFilledBuffer(x0, y0, x1, y1 int, c Cell) *Buffer {
	buf := NewBuffer()
	buf.Area.Min = image.Pt(x0, y0)
	buf.Area.Max = image.Pt(x1, y1)

	for x := buf.Area.Min.X; x < buf.Area.Max.X; x++ {
		for y := buf.Area.Min.Y; y < buf.Area.Max.Y; y++ {
			buf.SetCell(x, y, c)
		}
	}
	return buf
}

// Set assigns a char to (x,y)
func (b *Buffer) SetCell(x, y int, c Cell) {
	b.CellMap[image.Pt(x, y)] = c
}

func (b *Buffer) SetString(x, y int, s string, fg, bg Attribute) {
	for i, char := range s {
		b.SetCell(x+i, y, Cell{char, fg, bg})
	}
}

// At returns the cell at (x,y).
func (b *Buffer) At(x, y int) Cell {
	return b.CellMap[image.Pt(x, y)]
}

// Bounds returns the domain for which At can return non-zero color.
func (b *Buffer) Bounds() image.Rectangle {
	x0, y0, x1, y1 := 0, 0, 0, 0
	for p := range b.CellMap {
		if p.X > x1 {
			x1 = p.X
		}
		if p.X < x0 {
			x0 = p.X
		}
		if p.Y > y1 {
			y1 = p.Y
		}
		if p.Y < y0 {
			y0 = p.Y
		}
	}
	return image.Rect(x0, y0, x1+1, y1+1)
}

// SetArea assigns a new rect area to Buffer b.
func (b *Buffer) SetArea(r image.Rectangle) {
	b.Area.Max = r.Max
	b.Area.Min = r.Min
}

func (b *Buffer) SetAreaXY(x, y int) {
	b.Area.Min.Y = 0
	b.Area.Min.X = 0
	b.Area.Max.Y = y
	b.Area.Max.X = x
}

// Sync sets drawing area to the buffer's bound
func (b *Buffer) Sync() {
	b.SetArea(b.Bounds())
}

// Merge merges bs Buffers onto b
func (b *Buffer) Merge(bs ...*Buffer) {
	for _, buf := range bs {
		for p, c := range buf.CellMap {
			b.SetCell(p.X, p.Y, c)
		}
		b.SetArea(b.Area.Union(buf.Area))
	}
}

func (b *Buffer) MergeWithOffset(buf *Buffer, xOffset, yOffset int) {
	for p, c := range buf.CellMap {
		b.SetCell(p.X+xOffset, p.Y+yOffset, c)
	}
	rect := image.Rect(xOffset, yOffset, buf.Area.Max.X+xOffset, buf.Area.Max.Y+yOffset)
	b.SetArea(b.Area.Union(rect))
}

// Fill fills the Buffer b with ch,fg and bg.
func (b *Buffer) Fill(c Cell) {
	for x := b.Area.Min.X; x < b.Area.Max.X; x++ {
		for y := b.Area.Min.Y; y < b.Area.Max.Y; y++ {
			b.SetCell(x, y, c)
		}
	}
}
