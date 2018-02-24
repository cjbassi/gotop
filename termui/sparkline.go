package termui

var SPARKS = [8]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []int
	Title1     string
	Title2     string
	TitleColor Color
	Total      int
	LineColor  Color
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
		LineColor:  Theme.Sparkline,
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
		// prints sparkline
		for x := sl.X; x >= 1; x-- {
			char := SPARKS[0]
			if (sl.X - x) < len(line.Data) {
				char = SPARKS[int((float64(line.Data[(len(line.Data)-1)-(sl.X-x)])/float64(max))*7)]
			}
			buf.SetCell(x, sparkY, Cell{char, line.LineColor, sl.Bg})
		}
	}

	return buf
}
