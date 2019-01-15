// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import (
	"image"

	tb "github.com/nsf/termbox-go"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
}

func Render(items ...Drawable) {
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Draw(buf)
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Attrs.Fg)+1, tb.Attribute(cell.Attrs.Bg)+1,
				)
			}
		}
	}
	tb.Flush()
}

func Clear() {
	tb.Clear(tb.ColorDefault, tb.Attribute(Theme.Default.Bg+1))
}
