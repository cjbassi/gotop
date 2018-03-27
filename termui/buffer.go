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
func (self *Buffer) SetCell(x, y int, c Cell) {
	self.CellMap[image.Pt(x, y)] = c
}

// SetString assigns a string to a Buffer starting at (x,y).
func (self *Buffer) SetString(x, y int, s string, fg, bg Color) {
	for i, char := range s {
		self.SetCell(x+i, y, Cell{char, fg, bg})
	}
}

// At returns the cell at (x,y).
func (self *Buffer) At(x, y int) Cell {
	return self.CellMap[image.Pt(x, y)]
}

// SetArea assigns a new rect area to self.
func (self *Buffer) SetArea(r image.Rectangle) {
	self.Area.Max = r.Max
	self.Area.Min = r.Min
}

// SetAreaXY sets the Buffer bounds from (0,0) to (x,y).
func (self *Buffer) SetAreaXY(x, y int) {
	self.Area.Min.Y = 0
	self.Area.Min.X = 0
	self.Area.Max.Y = y
	self.Area.Max.X = x
}

// Merge merges the given buffers onto the current Buffer.
func (self *Buffer) Merge(bs ...*Buffer) {
	for _, buf := range bs {
		for p, c := range buf.CellMap {
			self.SetCell(p.X, p.Y, c)
		}
		self.SetArea(self.Area.Union(buf.Area))
	}
}

// MergeWithOffset merges a Buffer onto another with an offset.
func (self *Buffer) MergeWithOffset(buf *Buffer, xOffset, yOffset int) {
	for p, c := range buf.CellMap {
		self.SetCell(p.X+xOffset, p.Y+yOffset, c)
	}
	rect := image.Rect(xOffset, yOffset, buf.Area.Max.X+xOffset, buf.Area.Max.Y+yOffset)
	self.SetArea(self.Area.Union(rect))
}

// Fill fills the Buffer with a Cell.
func (self *Buffer) Fill(c Cell) {
	for x := self.Area.Min.X; x < self.Area.Max.X; x++ {
		for y := self.Area.Min.Y; y < self.Area.Max.Y; y++ {
			self.SetCell(x, y, c)
		}
	}
}
