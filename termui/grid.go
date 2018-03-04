package termui

var Body *Grid

// GridBufferer introduces a Bufferer that can be manipulated by Grid.
type GridBufferer interface {
	Bufferer
	Resize(int, int, int, int)
	SetGrid(int, int, int, int)
}

// Grid holds widgets and information about terminal dimensions.
// Widgets are adjusted and rendered through the grid.
type Grid struct {
	Widgets []GridBufferer
	Width   int
	Height  int
	Cols    int
	Rows    int
}

// NewGrid creates an empty Grid.
func NewGrid() *Grid {
	return &Grid{}
}

// Set assigns a widget and its grid dimensions to Grid.
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

// Resize resizes each widget in the grid.
func (g *Grid) Resize() {
	for _, w := range g.Widgets {
		w.Resize(g.Width, g.Height, g.Cols, g.Rows)
	}
}

// Buffer implements the Bufferer interface by merging each widget in Grid into one buffer.
func (g *Grid) Buffer() *Buffer {
	buf := NewFilledBuffer(0, 0, g.Width, g.Height, Cell{' ', ColorDefault, Theme.Bg})
	for _, w := range g.Widgets {
		buf.MergeWithOffset(w.Buffer(), w.GetXOffset(), w.GetYOffset())
	}
	return buf
}

// GetXOffset implements Bufferer interface.
func (g *Grid) GetXOffset() int {
	return 0
}

// GetYOffset implements Bufferer interface.
func (g *Grid) GetYOffset() int {
	return 0
}
