package termui

import (
	"image"
)

// Cell is a rune with assigned Fg and Bg.
type Cell struct {
	Ch rune
	Fg Color
	Bg Color
}

// Buffer is a renderable rectangle cell data container.
type Buffer struct {
	Area    image.Rectangle // selected drawing area
	CellMap map[image.Point]Cell
}

// NewCell returne a new Cell given all necessary fields.
func NewCell(ch rune, Fg, Bg Color) Cell {
	return Cell{ch, Fg, Bg}
}

// NewBuffer returns a new empty Buffer.
func NewBuffer() *Buffer {
	return &Buffer{
		CellMap: make(map[image.Point]Cell),
		Area:    image.Rectangle{},
	}
}

// NewFilledBuffer returns a new Buffer filled with the given Cell.
func NewFilledBuffer(x0, y0, x1, y1 int, c Cell) *Buffer {
	buf := NewBuffer()
	buf.Area.Min = image.Pt(x0, y0)
	buf.Area.Max = image.Pt(x1, y1)
	buf.Fill(c)
	return buf
}

// SetCell assigns a Cell to (x,y).
func (b *Buffer) SetCell(x, y int, c Cell) {
	b.CellMap[image.Pt(x, y)] = c
}

// SetString assigns a string to a Buffer starting at (x,y).
func (b *Buffer) SetString(x, y int, s string, fg, bg Color) {
	for i, char := range s {
		b.SetCell(x+i, y, Cell{char, fg, bg})
	}
}

// At returns the cell at (x,y).
func (b *Buffer) At(x, y int) Cell {
	return b.CellMap[image.Pt(x, y)]
}

// SetArea assigns a new rect area to Buffer b.
func (b *Buffer) SetArea(r image.Rectangle) {
	b.Area.Max = r.Max
	b.Area.Min = r.Min
}

// SetAreaXY sets the Buffer bounds from (0,0) to (x,y).
func (b *Buffer) SetAreaXY(x, y int) {
	b.Area.Min.Y = 0
	b.Area.Min.X = 0
	b.Area.Max.Y = y
	b.Area.Max.X = x
}

// Merge merges the given buffers onto the current Buffer.
func (b *Buffer) Merge(bs ...*Buffer) {
	for _, buf := range bs {
		for p, c := range buf.CellMap {
			b.SetCell(p.X, p.Y, c)
		}
		b.SetArea(b.Area.Union(buf.Area))
	}
}

// MergeWithOffset merges a Buffer onto another with an offset.
func (b *Buffer) MergeWithOffset(buf *Buffer, xOffset, yOffset int) {
	for p, c := range buf.CellMap {
		b.SetCell(p.X+xOffset, p.Y+yOffset, c)
	}
	rect := image.Rect(xOffset, yOffset, buf.Area.Max.X+xOffset, buf.Area.Max.Y+yOffset)
	b.SetArea(b.Area.Union(rect))
}

// Fill fills the Buffer with a Cell.
func (b *Buffer) Fill(c Cell) {
	for x := b.Area.Min.X; x < b.Area.Max.X; x++ {
		for y := b.Area.Min.Y; y < b.Area.Max.Y; y++ {
			b.SetCell(x, y, c)
		}
	}
}
