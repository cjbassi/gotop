package termui

import "fmt"

// BarChart creates multiple bars in a widget:
type List struct {
	*Block
	TextColor  Attribute
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
		bg := ColorGreen
		if bc.Data[y] >= bc.Threshold {
			bg = ColorRed
		}
		r := MaxString(text, (bc.X - 4))
		buf.SetString(1, y+1, r, ColorWhite, ColorDefault)
		buf.SetString(bc.X-2, y+1, fmt.Sprintf("%dC", bc.Data[y]), bg, ColorDefault)
	}

	return buf
}
