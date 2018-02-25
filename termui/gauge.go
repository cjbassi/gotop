package termui

import (
	"strconv"
)

// Gauge is a progress bar like widget.
type Gauge struct {
	*Block
	Percent     int
	GaugeColor  Color
	Description string
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	return &Gauge{
		Block:      NewBlock(),
		GaugeColor: Theme.GaugeColor,
	}
}

// Buffer implements Bufferer interface.
func (g *Gauge) Buffer() *Buffer {
	buf := g.Block.Buffer()

	// plot bar
	width := g.Percent * g.X / 100
	for y := 1; y <= g.Y; y++ {
		for x := 1; x <= width; x++ {
			buf.SetCell(x, y, Cell{' ', g.GaugeColor, g.GaugeColor})
		}
	}

	// plot percentage
	s := strconv.Itoa(g.Percent) + "%" + g.Description
	s = MaxString(s, g.X)

	y := (g.Y + 1) / 2
	x := ((g.X - len(s)) + 1) / 2

	for i, char := range s {
		bg := g.Bg
		fg := g.Fg
		if x+i < width {
			fg = g.GaugeColor
			bg = AttrReverse
		}
		buf.SetCell(1+x+i, y, Cell{char, fg, bg})
	}

	return buf
}
