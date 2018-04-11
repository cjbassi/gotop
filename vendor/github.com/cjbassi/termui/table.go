package termui

import (
	"fmt"
	"strings"
)

// Table tracks all the attributes of a Table instance
type Table struct {
	*Block
	Header       []string
	Rows         [][]string
	ColWidths    []int
	CellXPos     []int // column position
	Gap          int   // gap between columns
	Cursor       Color
	UniqueCol    int    // the column used to identify the selected item
	SelectedItem string // used to keep the cursor on the correct item if the data changes
	SelectedRow  int
	TopRow       int    // used to indicate where in the table we are scrolled at
	ColResizer   func() // for widgets that inherit a Table and want to overload the ColResize method
}

// NewTable returns a new Table instance
func NewTable() *Table {
	self := &Table{
		Block:       NewBlock(),
		Cursor:      Theme.TableCursor,
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
	// calculate gap size based on total width
	self.Gap = 3
	if self.X < 50 {
		self.Gap = 1
	} else if self.X < 75 {
		self.Gap = 2
	}

	cur := 0
	for _, w := range self.ColWidths {
		cur += self.Gap
		self.CellXPos = append(self.CellXPos, cur)
		cur += w
	}
}

// Buffer implements the Bufferer interface.
func (self *Table) Buffer() *Buffer {
	buf := self.Block.Buffer()

	// removes gap at the bottom of the current view if there is one
	// if self.TopRow > len(self.Rows)-(self.Y-1) {
	// 	self.TopRow = len(self.Rows) - (self.Y - 1)
	// }

	self.ColResizer()

	// prints header
	for i, width := range self.ColWidths {
		if width == 0 {
			break
		}
		r := MaxString(self.Header[i], self.X-6)
		buf.SetString(self.CellXPos[i], 1, r, self.Fg|AttrBold, self.Bg)
	}

	// prints each row
	for rowNum := self.TopRow; rowNum < self.TopRow+self.Y-1 && rowNum < len(self.Rows); rowNum++ {
		if rowNum < 0 || rowNum >= len(self.Rows) {
			Error("table rows",
				fmt.Sprint(
					"rowNum: ", rowNum, "\n",
					"self.TopRow: ", self.TopRow, "\n",
					"len(self.Rows): ", len(self.Rows), "\n",
					"self.Y: ", self.Y,
				))
		}
		row := self.Rows[rowNum]
		y := (rowNum + 2) - self.TopRow

		// prints cursor
		bg := self.Bg
		if (self.SelectedItem == "" && rowNum == self.SelectedRow) || (self.SelectedItem != "" && self.SelectedItem == row[self.UniqueCol]) {
			bg = self.Cursor
			for _, width := range self.ColWidths {
				if width == 0 {
					break
				}
				buf.SetString(1, y, strings.Repeat(" ", self.X), self.Fg, bg)
			}
			self.SelectedItem = row[self.UniqueCol]
			self.SelectedRow = rowNum
		}

		// prints each col of the row
		for i, width := range self.ColWidths {
			if width == 0 {
				break
			}
			r := MaxString(row[i], self.X-6)
			buf.SetString(self.CellXPos[i], y, r, self.Fg, bg)
		}
	}

	return buf
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
	if self.SelectedRow > self.TopRow+(self.Y-2) {
		self.TopRow = self.SelectedRow - (self.Y - 2)
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
	self.SelectedRow = self.SelectedRow - (self.Y-2)/2
	self.calcPos()
}

func (self *Table) HalfPageDown() {
	self.SelectedRow = self.SelectedRow + (self.Y-2)/2
	self.calcPos()
}

func (self *Table) PageUp() {
	self.SelectedRow -= (self.Y - 2)
	self.calcPos()
}

func (self *Table) PageDown() {
	self.SelectedRow += (self.Y - 2)
	self.calcPos()
}

func (self *Table) Click(x, y int) {
	x = x - self.XOffset
	y = y - self.YOffset
	if (x > 0 && x <= self.X) && (y > 0 && y <= self.Y) {
		self.SelectedRow = (self.TopRow + y) - 2
		self.calcPos()
	}
}
