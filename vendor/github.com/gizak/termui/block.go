// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"
)

type Block struct {
	Border       bool
	BorderAttrs  AttrPair
	BorderLeft   bool
	BorderRight  bool
	BorderTop    bool
	BorderBottom bool

	image.Rectangle
	Inner image.Rectangle

	Title      string
	TitleAttrs AttrPair
}

// NewBlock returns a *Block which inherits styles from current theme.
func NewBlock() *Block {
	return &Block{
		Border:       true,
		BorderAttrs:  Theme.Block.Border,
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		TitleAttrs: Theme.Block.Title,
	}
}

func (b *Block) drawBorder(buf *Buffer) {
	if !b.Border {
		return
	}

	verticalCell := Cell{VERTICAL_LINE, b.BorderAttrs}
	horizontalCell := Cell{HORIZONTAL_LINE, b.BorderAttrs}

	// draw lines
	if b.BorderTop {
		buf.Fill(horizontalCell, image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Min.Y+1))
	}
	if b.BorderBottom {
		buf.Fill(horizontalCell, image.Rect(b.Min.X, b.Max.Y-1, b.Max.X, b.Max.Y))
	}
	if b.BorderLeft {
		buf.Fill(verticalCell, image.Rect(b.Min.X, b.Min.Y, b.Min.X+1, b.Max.Y))
	}
	if b.BorderRight {
		buf.Fill(verticalCell, image.Rect(b.Max.X-1, b.Min.Y, b.Max.X, b.Max.Y))
	}

	// draw corners
	if b.BorderTop && b.BorderLeft {
		buf.SetCell(Cell{TOP_LEFT, b.BorderAttrs}, b.Min)
	}
	if b.BorderTop && b.BorderRight {
		buf.SetCell(Cell{TOP_RIGHT, b.BorderAttrs}, image.Pt(b.Max.X-1, b.Min.Y))
	}
	if b.BorderBottom && b.BorderLeft {
		buf.SetCell(Cell{BOTTOM_LEFT, b.BorderAttrs}, image.Pt(b.Min.X, b.Max.Y-1))
	}
	if b.BorderBottom && b.BorderRight {
		buf.SetCell(Cell{BOTTOM_RIGHT, b.BorderAttrs}, b.Max.Sub(image.Pt(1, 1)))
	}
}

func (b *Block) Draw(buf *Buffer) {
	b.drawBorder(buf)
	buf.SetString(
		b.Title,
		b.TitleAttrs,
		image.Pt(b.Min.X+2, b.Min.Y),
	)
}

func (b *Block) SetRect(x1, y1, x2, y2 int) {
	b.Rectangle = image.Rect(x1, y1, x2, y2)
	b.Inner = image.Rect(b.Min.X+1, b.Min.Y+1, b.Max.X-1, b.Max.Y-1)
}

func (b *Block) GetRect() image.Rectangle {
	return b.Rectangle
}
