package termui

import (
	"image"
)

// Block is a base struct for all other upper level widgets.
type Block struct {
	Grid     image.Rectangle
	X        int // largest X value in the inner square
	Y        int // largest Y value in the inner square
	XOffset  int // the X position of the widget on the terminal
	YOffset  int // the Y position of the widget on the terminal
	Label    string
	BorderFg Color
	BorderBg Color
	LabelFg  Color
	LabelBg  Color
	Fg       Color
	Bg       Color
}

// NewBlock returns a *Block which inherits styles from the current theme.
func NewBlock() *Block {
	return &Block{
		Fg:       Theme.Fg,
		Bg:       Theme.Bg,
		BorderFg: Theme.BorderFg,
		BorderBg: Theme.BorderBg,
		LabelFg:  Theme.LabelFg,
		LabelBg:  Theme.LabelBg,
	}
}

func (self *Block) drawBorder(buf *Buffer) {
	x := self.X + 1
	y := self.Y + 1

	// draw lines
	buf.Merge(NewFilledBuffer(0, 0, x, 1, Cell{HORIZONTAL_LINE, self.BorderFg, self.BorderBg}))
	buf.Merge(NewFilledBuffer(0, y, x, y+1, Cell{HORIZONTAL_LINE, self.BorderFg, self.BorderBg}))
	buf.Merge(NewFilledBuffer(0, 0, 1, y+1, Cell{VERTICAL_LINE, self.BorderFg, self.BorderBg}))
	buf.Merge(NewFilledBuffer(x, 0, x+1, y+1, Cell{VERTICAL_LINE, self.BorderFg, self.BorderBg}))

	// draw corners
	buf.SetCell(0, 0, Cell{TOP_LEFT, self.BorderFg, self.BorderBg})
	buf.SetCell(x, 0, Cell{TOP_RIGHT, self.BorderFg, self.BorderBg})
	buf.SetCell(0, y, Cell{BOTTOM_LEFT, self.BorderFg, self.BorderBg})
	buf.SetCell(x, y, Cell{BOTTOM_RIGHT, self.BorderFg, self.BorderBg})
}

func (self *Block) drawLabel(buf *Buffer) {
	r := MaxString(self.Label, (self.X-3)-1)
	buf.SetString(3, 0, r, self.LabelFg, self.LabelBg)
	if self.Label == "" {
		return
	}
	c := Cell{' ', self.Fg, self.Bg}
	buf.SetCell(2, 0, c)
	if len(self.Label)+3 < self.X {
		buf.SetCell(len(self.Label)+3, 0, c)
	} else {
		buf.SetCell(self.X-1, 0, c)
	}
}

// Resize computes Height, Width, XOffset, and YOffset given terminal dimensions.
func (self *Block) Resize(termWidth, termHeight, termCols, termRows int) {
	self.X = int((float64(self.Grid.Dx())/float64(termCols))*float64(termWidth)) - 2
	self.Y = int((float64(self.Grid.Dy())/float64(termRows))*float64(termHeight)) - 2
	self.XOffset = int((float64(self.Grid.Min.X) / float64(termCols)) * float64(termWidth))
	self.YOffset = int((float64(self.Grid.Min.Y) / float64(termRows)) * float64(termHeight))
}

// SetGrid create a rectangle representing the block's dimensions in the grid.
func (self *Block) SetGrid(c0, r0, c1, r1 int) {
	self.Grid = image.Rect(c0, r0, c1, r1)
}

// GetXOffset implements Bufferer interface.
func (self *Block) GetXOffset() int {
	return self.XOffset
}

// GetYOffset implements Bufferer interface.
func (self *Block) GetYOffset() int {
	return self.YOffset
}

// Buffer implements Bufferer interface and draws background, border, and borderlabel.
func (self *Block) Buffer() *Buffer {
	buf := NewBuffer()
	buf.SetAreaXY(self.X+2, self.Y+2)
	buf.Fill(Cell{' ', ColorDefault, self.Bg})

	self.drawBorder(buf)
	self.drawLabel(buf)

	return buf
}
