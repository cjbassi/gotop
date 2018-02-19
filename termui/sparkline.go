package termui

import (
	"fmt"
	"github.com/cjbassi/gotop/utils"
)

var SPARKS = [8]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []int
	Title      string
	TitleColor Attribute
	Total      int
	LineColor  Attribute
}

// Sparklines is a renderable widget which groups together the given sparklines.
type Sparklines struct {
	*Block
	Lines []*Sparkline
}

// Add appends a given Sparkline to s *Sparklines.
func (s *Sparklines) Add(sl Sparkline) {
	s.Lines = append(s.Lines, &sl)
}

// NewSparkline returns a unrenderable single sparkline that intended to be added into Sparklines.
func NewSparkline() *Sparkline {
	return &Sparkline{
		TitleColor: Theme.Fg,
		LineColor:  Theme.SparkLine,
	}
}

// NewSparklines return a new *Sparklines with given Sparkline(s), you can always add a new Sparkline later.
func NewSparklines(ss ...*Sparkline) *Sparklines {
	return &Sparklines{
		Block: NewBlock(),
		Lines: ss,
	}
}

// Buffer implements Bufferer interface.
func (sl *Sparklines) Buffer() *Buffer {
	buf := sl.Block.Buffer()

	lc := len(sl.Lines) // lineCount

	// for each line
	for i, line := range sl.Lines {

		// Total and current
		y := 2 + (sl.Y/lc)*i
		total := ""
		title := ""
		current := ""

		cur := line.Data[len(line.Data)-1]
		curMag := "B"
		if cur >= 1000000 {
			cur = int(utils.BytesToMB(cur))
			curMag = "MB"
		} else if cur >= 1000 {
			cur = int(utils.BytesToKB(cur))
			curMag = "kB"
		}

		t := line.Total
		tMag := "B"
		if t >= 1000000000 {
			t = int(utils.BytesToGB(t))
			tMag = "GB"
		} else if t >= 1000000 {
			t = int(utils.BytesToMB(t))
			tMag = "MB"
		}

		if i == 0 {
			total = fmt.Sprintf(" Total Rx: %3d %s", t, tMag)
			current = fmt.Sprintf(" Rx/s: %7d %2s/s", cur, curMag)
		} else {
			total = fmt.Sprintf(" Total Tx: %3d %s", t, tMag)
			current = fmt.Sprintf(" Tx/s: %7d %2s/s", cur, curMag)
		}

		total = MaxString(total, sl.X)
		title = MaxString(current, sl.X)
		buf.SetString(1, y, total, line.TitleColor|AttrBold, sl.Bg)
		buf.SetString(1, y+1, title, line.TitleColor|AttrBold, sl.Bg)

		// sparkline
		y = (sl.Y / lc) * (i + 1)
		// finds max used for relative heights
		max := 1
		for i := len(line.Data) - 1; i >= 0 && sl.X-((len(line.Data)-1)-i) >= 1; i-- {
			if line.Data[i] > max {
				max = line.Data[i]
			}
		}
		// prints sparkline
		for x := sl.X; x >= 1; x-- {
			char := SPARKS[0]
			if (sl.X - x) < len(line.Data) {
				char = SPARKS[int((float64(line.Data[(len(line.Data)-1)-(sl.X-x)])/float64(max))*7)]
			}
			buf.SetCell(x, y, Cell{char, line.LineColor, sl.Bg})
		}
	}

	return buf
}
