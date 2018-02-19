package termui

import "strconv"

// Gauge is a progress bar like widget.
type Gauge struct {
	*Block
	Percent      int
	BarColor     Attribute
	PercentColor Attribute
	Description  string
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	return &Gauge{
		Block:        NewBlock(),
		PercentColor: Theme.Fg,
		BarColor:     Theme.Bg,
	}
}

// Buffer implements Bufferer interface.
func (g *Gauge) Buffer() *Buffer {
	buf := g.Block.Buffer()

	// plot bar
	width := g.Percent * g.X / 100
	for y := 1; y <= g.Y; y++ {
		for x := 1; x <= width; x++ {
			bg := g.BarColor
			if bg == ColorDefault {
				bg |= AttrReverse
			}
			buf.SetCell(x, y, Cell{' ', ColorDefault, bg})
		}
	}

	// plot percentage
	s := strconv.Itoa(g.Percent) + "%" + g.Description
	y := (g.Y + 1) / 2
	s = MaxString(s, g.X)
	x := ((g.X - len(s)) + 1) / 2

	for i, char := range s {
		bg := g.Bg
		if x+i < width {
			bg = AttrReverse
		}
		buf.SetCell(1+x+i, y, Cell{char, g.PercentColor, bg})
	}

	return buf
}
