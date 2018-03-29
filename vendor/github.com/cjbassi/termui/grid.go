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
func (self *Grid) Set(x0, y0, x1, y1 int, widget GridBufferer) {
	if widget == nil {
		return
	}
	if x1 <= x0 || y1 <= y0 {
		panic("Invalid widget coordinates")
	}

	widget.SetGrid(x0, y0, x1, y1)
	widget.Resize(self.Width, self.Height, self.Cols, self.Rows)

	self.Widgets = append(self.Widgets, widget)
}

// Resize resizes each widget in the grid.
func (self *Grid) Resize() {
	for _, w := range self.Widgets {
		w.Resize(self.Width, self.Height, self.Cols, self.Rows)
	}
}

// Buffer implements the Bufferer interface by merging each widget in Grid into one buffer.
func (self *Grid) Buffer() *Buffer {
	buf := NewFilledBuffer(0, 0, self.Width, self.Height, Cell{' ', ColorDefault, Theme.Bg})
	for _, w := range self.Widgets {
		buf.MergeWithOffset(w.Buffer(), w.GetXOffset(), w.GetYOffset())
	}
	return buf
}

// GetXOffset implements Bufferer interface.
func (self *Grid) GetXOffset() int {
	return 0
}

// GetYOffset implements Bufferer interface.
func (self *Grid) GetYOffset() int {
	return 0
}
