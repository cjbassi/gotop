package termui

var Body *Grid

// GridBufferer introduces a Bufferer that can be manipulated by Grid.
type GridBufferer interface {
	Bufferer
	Resize(int, int, int, int)
	SetGrid(int, int, int, int)
}

type Grid struct {
	Widgets []GridBufferer
	Width   int
	Height  int
	Cols    int
	Rows    int
	BgColor Color
}

func NewGrid() *Grid {
	return &Grid{}
}

func (g *Grid) Set(x0, y0, x1, y1 int, widget GridBufferer) {
	if widget == nil {
		return
	}
	if x1 <= x0 || y1 <= y0 {
		panic("Invalid widget coordinates")
	}

	widget.SetGrid(x0, y0, x1, y1)
	widget.Resize(g.Width, g.Height, g.Cols, g.Rows)

	g.Widgets = append(g.Widgets, widget)
}

func (g *Grid) Resize() {
	for _, w := range g.Widgets {
		w.Resize(g.Width, g.Height, g.Cols, g.Rows)
	}
}

// Buffer implements Bufferer interface.
func (g *Grid) Buffer() *Buffer {
	buf := NewBuffer()
	for _, w := range g.Widgets {
		buf.MergeWithOffset(w.Buffer(), w.GetXOffset(), w.GetYOffset())
	}
	return buf
}

func (g *Grid) GetXOffset() int {
	return 0
}

func (g *Grid) GetYOffset() int {
	return 0
}
