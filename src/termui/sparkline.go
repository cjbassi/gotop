package termui

import (
	"fmt"
)

var SPARKS = [8]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []int
	Title1     string
	Title2     string
	TitleColor Color
	LineColor  Color
}

// Sparklines is a renderable widget which groups together the given sparklines.
type Sparklines struct {
	*Block
	Lines []*Sparkline
}

// Add appends a given Sparkline to the *Sparklines.
func (self *Sparklines) Add(sl Sparkline) {
	self.Lines = append(self.Lines, &sl)
}

// NewSparkline returns an unrenderable single sparkline that intended to be added into a Sparklines.
func NewSparkline() *Sparkline {
	return &Sparkline{
		TitleColor: Theme.Fg,
		LineColor:  Theme.Sparkline,
	}
}

// NewSparklines return a new *Sparklines with given Sparklines, you can always add a new Sparkline later.
func NewSparklines(ss ...*Sparkline) *Sparklines {
	return &Sparklines{
		Block: NewBlock(),
		Lines: ss,
	}
}

// Buffer implements Bufferer interface.
func (self *Sparklines) Buffer() *Buffer {
	buf := self.Block.Buffer()

	lc := len(self.Lines) // lineCount

	// renders each sparkline and its titles
	for i, line := range self.Lines {

		// prints titles
		title1Y := 2 + (self.Y/lc)*i
		title2Y := (2 + (self.Y/lc)*i) + 1
		title1 := MaxString(line.Title1, self.X)
		title2 := MaxString(line.Title2, self.X)
		buf.SetString(1, title1Y, title1, line.TitleColor|AttrBold, self.Bg)
		buf.SetString(1, title2Y, title2, line.TitleColor|AttrBold, self.Bg)

		sparkY := (self.Y / lc) * (i + 1)
		// finds max data in current view used for relative heights
		max := 1
		for i := len(line.Data) - 1; i >= 0 && self.X-((len(line.Data)-1)-i) >= 1; i-- {
			if line.Data[i] > max {
				max = line.Data[i]
			}
		}
		// prints sparkline
		for x := self.X; x >= 1; x-- {
			char := SPARKS[0]
			if (self.X - x) < len(line.Data) {
				offset := self.X - x
				cur_item := line.Data[(len(line.Data)-1)-offset]
				percent := float64(cur_item) / float64(max)
				index := int(percent * 7)
				if index < 0 || index >= len(SPARKS) {
					Error("sparkline",
						fmt.Sprint(
							"len(line.Data): ", len(line.Data), "\n",
							"max: ", max, "\n",
							"x: ", x, "\n",
							"self.X: ", self.X, "\n",
							"offset: ", offset, "\n",
							"cur_item: ", cur_item, "\n",
							"percent: ", percent, "\n",
							"index: ", index, "\n",
							"len(SPARKS): ", len(SPARKS),
						))
				}
				char = SPARKS[index]
			}
			buf.SetCell(x, sparkY, Cell{char, line.LineColor, self.Bg})
		}
	}

	return buf
}
