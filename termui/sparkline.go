package termui

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
func (s *Sparklines) Add(sl Sparkline) {
	s.Lines = append(s.Lines, &sl)
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
func (sl *Sparklines) Buffer() *Buffer {
	buf := sl.Block.Buffer()

	lc := len(sl.Lines) // lineCount

	// renders each sparkline and its titles
	for i, line := range sl.Lines {

		// prints titles
		title1Y := 2 + (sl.Y/lc)*i
		title2Y := (2 + (sl.Y/lc)*i) + 1
		title1 := MaxString(line.Title1, sl.X)
		title2 := MaxString(line.Title2, sl.X)
		buf.SetString(1, title1Y, title1, line.TitleColor|AttrBold, sl.Bg)
		buf.SetString(1, title2Y, title2, line.TitleColor|AttrBold, sl.Bg)

		sparkY := (sl.Y / lc) * (i + 1)
		// finds max data in current view used for relative heights
		max := 1
		for i := len(line.Data) - 1; i >= 0 && sl.X-((len(line.Data)-1)-i) >= 1; i-- {
			if line.Data[i] > max {
				max = line.Data[i]
			}
		}
		maxHeight := sparkY - title2Y
		gap := 100 / float64(maxHeight)
		// prints sparkline
		for x := sl.X; x >= 1; x-- {
			var char rune
			if (sl.X - x) < len(line.Data) {
				cur := line.Data[(len(line.Data)-1)-(sl.X-x)]
				percent := (float64(cur) / float64(max)) * 100
				for i := 0; i < maxHeight; i++ {
					min := float64(i) * gap
					y := sparkY - i
					buf.SetCell(x, y, Cell{char, line.LineColor, sl.Bg})
					if percent < min {
						char = ' '
					} else if percent > min+gap {
						char = SPARKS[7]
					} else {
						char = SPARKS[int(((percent-min)/gap)*7)]
					}
					buf.SetCell(x, y, Cell{char, line.LineColor, sl.Bg})
				}
			} else {
				char = SPARKS[0]
				buf.SetCell(x, sparkY, Cell{char, line.LineColor, sl.Bg})
			}
		}
	}

	return buf
}
