package termui

import (
	"strconv"

	"github.com/gdamore/tcell"
)

// Gauge is a progress bar like widget.
type Gauge struct {
	*Block
	Percent     int
	BarStyle    tcell.Style
	Description string
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	return &Gauge{
		Block:    NewBlock(),
		BarStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Reverse(true),
	}
}

// Buffer implements Bufferer interface.
func (self *Gauge) Buffer() *Buffer {
	buf := self.Block.Buffer()

	// plot bar
	width := self.Percent * self.X / 100
	for y := 1; y <= self.Y; y++ {
		for x := 1; x <= width; x++ {
			buf.SetCell(x, y, Cell{' ', self.BarStyle})
		}
	}

	// plot percentage
	s := strconv.Itoa(self.Percent) + "%" + self.Description
	s = MaxString(s, self.X)
	y := (self.Y + 1) / 2
	x := ((self.X - len(s)) + 1) / 2
	for i, char := range s {
		st := tcell.StyleDefault
		if x+i < width {
			st = st.Reverse(true)
		}
		buf.SetCell(1+x+i, y, Cell{char, st})
	}

	return buf
}
