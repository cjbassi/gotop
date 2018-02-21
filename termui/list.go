package termui

import (
	"fmt"
)

// BarChart creates multiple bars in a widget:
type List struct {
	*Block
	TextColor  Color
	Data       []int
	DataLabels []string
	Threshold  int
}

// NewBarChart returns a new *BarChart with current theme.
func NewList() *List {
	return &List{
		Block:     NewBlock(),
		TextColor: Theme.Fg,
	}
}

// Buffer implements Bufferer interface.
func (bc *List) Buffer() *Buffer {
	buf := bc.Block.Buffer()

	for y, text := range bc.DataLabels {
		if y+1 > bc.Y {
			break
		}
		fg := Theme.TempLow
		if bc.Data[y] >= bc.Threshold {
			fg = Theme.TempHigh
		}
		r := MaxString(text, (bc.X - 4))
		buf.SetString(1, y+1, r, Theme.Fg, bc.Bg)
		buf.SetString(bc.X-2, y+1, fmt.Sprintf("%dC", bc.Data[y]), fg, bc.Bg)
	}

	return buf
}
