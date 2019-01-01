package termui

import (
	"image"
	"strings"
	"sync"

	. "github.com/gizak/termui"
)

// Table tracks all the attributes of a Table instance
type Table struct {
	*Block

	sync.Mutex

	Header []string
	Rows   [][]string

	ColWidths  []int
	CellXPos   []int  // column position
	ColResizer func() // for widgets that inherit a Table and want to overload the ColResize method
	Gap        int    // gap between columns
	PadLeft    int

	Cursor      bool
	CursorColor Attribute

	UniqueCol    int    // the column used to identify the selected item
	SelectedItem string // used to keep the cursor on the correct item if the data changes
	SelectedRow  int
	TopRow       int // used to indicate where in the table we are scrolled at
}

// NewTable returns a new Table instance
func NewTable() *Table {
	self := &Table{
		Block: NewBlock(),
		// CursorColor: Theme.TableCursor,
		SelectedRow: 0,
		TopRow:      0,
		UniqueCol:   0,
	}
	self.ColResizer = self.ColResize
	return self
}

// ColResize is the default column resizer, but can be overriden.
// ColResize calculates the width of each column.
func (self *Table) ColResize() {
}

func (self *Table) Draw(buf *Buffer) {
	self.Lock()

	self.Block.Draw(buf)

	self.ColResizer()

	// finds exact column starting position
	self.CellXPos = []int{}
	cur := 1 + self.PadLeft
	for _, w := range self.ColWidths {
		self.CellXPos = append(self.CellXPos, cur)
		cur += w
		cur += self.Gap
	}

	// prints header
	for i, h := range self.Header {
		width := self.ColWidths[i]
		if width == 0 {
			continue
		}
		// don't render column if it doesn't fit in widget
		if width > (self.Inner.Dx()-self.CellXPos[i])+1 {
			continue
		}
		buf.SetString(
			h,
			image.Pt(self.Inner.Min.X+self.CellXPos[i]-1, self.Inner.Min.Y),
			AttrPair{Theme.Default.Fg | AttrBold, -1},
		)
	}

	// prints each row
	for rowNum := self.TopRow; rowNum < self.TopRow+self.Inner.Dy()-1 && rowNum < len(self.Rows); rowNum++ {
		if rowNum < 0 || rowNum >= len(self.Rows) {
			panic("TODO")
		}
		row := self.Rows[rowNum]
		y := (rowNum + 2) - self.TopRow

		// prints cursor
		fg := Theme.Default.Fg
		if self.Cursor {
			if (self.SelectedItem == "" && rowNum == self.SelectedRow) || (self.SelectedItem != "" && self.SelectedItem == row[self.UniqueCol]) {
				fg = self.CursorColor | AttrReverse
				for _, width := range self.ColWidths {
					if width == 0 {
						continue
					}
					buf.SetString(
						strings.Repeat(" ", self.Inner.Dx()),
						image.Pt(self.Inner.Min.X, self.Inner.Min.Y+y-1),
						AttrPair{fg, -1},
					)
				}
				self.SelectedItem = row[self.UniqueCol]
				self.SelectedRow = rowNum
			}
		}

		// prints each col of the row
		for i, width := range self.ColWidths {
			if width == 0 {
				continue
			}
			// don't render column if width is greater than distance to end of widget
			if width > (self.Inner.Dx()-self.CellXPos[i])+1 {
				continue
			}
			r := TrimString(row[i], width)
			buf.SetString(
				r,
				image.Pt(self.Inner.Min.X+self.CellXPos[i]-1, self.Inner.Min.Y+y-1),
				AttrPair{fg, -1},
			)
		}
	}
	self.Unlock()
}

/////////////////////////////////////////////////////////////////////////////////
//                               Cursor Movement                               //
/////////////////////////////////////////////////////////////////////////////////

// calcPos is used to calculate the cursor position and the current view.
func (self *Table) calcPos() {
	self.SelectedItem = ""

	if self.SelectedRow < 0 {
		self.SelectedRow = 0
	}
	if self.SelectedRow < self.TopRow {
		self.TopRow = self.SelectedRow
	}

	if self.SelectedRow > len(self.Rows)-1 {
		self.SelectedRow = len(self.Rows) - 1
	}
	if self.SelectedRow > self.TopRow+(self.Inner.Dy()-2) {
		self.TopRow = self.SelectedRow - (self.Inner.Dy() - 2)
	}
}

func (self *Table) Up() {
	self.SelectedRow -= 1
	self.calcPos()
}

func (self *Table) Down() {
	self.SelectedRow += 1
	self.calcPos()
}

func (self *Table) Top() {
	self.SelectedRow = 0
	self.calcPos()
}

func (self *Table) Bottom() {
	self.SelectedRow = len(self.Rows) - 1
	self.calcPos()
}

// The number of lines in a page is equal to the height of the widgeself.

func (self *Table) HalfPageUp() {
	self.SelectedRow = self.SelectedRow - (self.Inner.Dy()-2)/2
	self.calcPos()
}

func (self *Table) HalfPageDown() {
	self.SelectedRow = self.SelectedRow + (self.Inner.Dy()-2)/2
	self.calcPos()
}

func (self *Table) PageUp() {
	self.SelectedRow -= (self.Inner.Dy() - 2)
	self.calcPos()
}

func (self *Table) PageDown() {
	self.SelectedRow += (self.Inner.Dy() - 2)
	self.calcPos()
}

func (self *Table) Click(x, y int) {
	x = x - self.Min.X
	y = y - self.Min.Y
	if (x > 0 && x <= self.Inner.Dx()) && (y > 0 && y <= self.Inner.Dy()) {
		self.SelectedRow = (self.TopRow + y) - 2
		self.calcPos()
	}
}
