// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
)

// Cell represents a terminal cell and is a rune with Fg and Bg Attributes
type Cell struct {
	Rune  rune
	Attrs AttrPair
}

// Buffer represents a section of a terminal and is a renderable rectangle cell data container.
type Buffer struct {
	image.Rectangle
	CellMap map[image.Point]Cell
}

func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,
		CellMap:   make(map[image.Point]Cell),
	}
	buf.Fill(Cell{' ', AttrPair{ColorDefault, ColorDefault}}, r) // clears out area
	return buf
}

func (b *Buffer) GetCell(p image.Point) Cell {
	return b.CellMap[p]
}

func (b *Buffer) SetCell(c Cell, p image.Point) {
	b.CellMap[p] = c
}

func (b *Buffer) Fill(c Cell, rect image.Rectangle) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			b.SetCell(c, image.Pt(x, y))
		}
	}
}

func (b *Buffer) SetString(s string, pair AttrPair, p image.Point) {
	for i, char := range s {
		b.SetCell(Cell{char, pair}, image.Pt(p.X+i, p.Y))
	}
}
