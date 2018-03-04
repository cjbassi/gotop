package termui

import (
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
	t := &Table{
		Block:       NewBlock(),
		Cursor:      Theme.TableCursor,
		SelectedRow: 0,
		TopRow:      0,
		UniqueCol:   0,
	}
	t.ColResizer = t.ColResize
	return t
}

// ColResize is the default column resizer, but can be overriden.
// ColResize calculates the width of each column.
func (t *Table) ColResize() {
	// calculate gap size based on total width
	t.Gap = 3
	if t.X < 50 {
		t.Gap = 1
	} else if t.X < 75 {
		t.Gap = 2
	}

	cur := 0
	for _, w := range t.ColWidths {
		cur += t.Gap
		t.CellXPos = append(t.CellXPos, cur)
		cur += w
	}
}

// Buffer implements the Bufferer interface.
func (t *Table) Buffer() *Buffer {
	buf := t.Block.Buffer()

	// removes gap at the bottom of the current view if there is one
	if t.TopRow > len(t.Rows)-(t.Y-1) {
		t.TopRow = len(t.Rows) - (t.Y - 1)
	}

	t.ColResizer()

	// prints header
	for i, width := range t.ColWidths {
		if width == 0 {
			break
		}
		r := MaxString(t.Header[i], t.X-6)
		buf.SetString(t.CellXPos[i], 1, r, t.Fg|AttrBold, t.Bg)
	}

	// prints each row
	for rowNum := t.TopRow; rowNum < t.TopRow+t.Y-1 && rowNum < len(t.Rows); rowNum++ {
		row := t.Rows[rowNum]
		y := (rowNum + 2) - t.TopRow

		// prints cursor
		bg := t.Bg
		if (t.SelectedItem == "" && rowNum == t.SelectedRow) || (t.SelectedItem != "" && t.SelectedItem == row[t.UniqueCol]) {
			bg = t.Cursor
			for _, width := range t.ColWidths {
				if width == 0 {
					break
				}
				buf.SetString(1, y, strings.Repeat(" ", t.X), t.Fg, bg)
			}
			t.SelectedItem = row[t.UniqueCol]
			t.SelectedRow = rowNum
		}

		// prints each col of the row
		for i, width := range t.ColWidths {
			if width == 0 {
				break
			}
			r := MaxString(row[i], t.X-6)
			buf.SetString(t.CellXPos[i], y, r, t.Fg, bg)
		}
	}

	return buf
}

/////////////////////////////////////////////////////////////////////////////////
//                               Cursor Movement                               //
/////////////////////////////////////////////////////////////////////////////////

// calcPos is used to calculate the cursor position and the current view.
func (t *Table) calcPos() {
	t.SelectedItem = ""

	if t.SelectedRow < 0 {
		t.SelectedRow = 0
	}
	if t.SelectedRow < t.TopRow {
		t.TopRow = t.SelectedRow
	}

	if t.SelectedRow > len(t.Rows)-1 {
		t.SelectedRow = len(t.Rows) - 1
	}
	if t.SelectedRow > t.TopRow+(t.Y-2) {
		t.TopRow = t.SelectedRow - (t.Y - 2)
	}
}

func (t *Table) Up() {
	t.SelectedRow -= 1
	t.calcPos()
}

func (t *Table) Down() {
	t.SelectedRow += 1
	t.calcPos()
}

func (t *Table) Top() {
	t.SelectedRow = 0
	t.calcPos()
}

func (t *Table) Bottom() {
	t.SelectedRow = len(t.Rows) - 1
	t.calcPos()
}

// The number of lines in a page is equal to the height of the widget.

func (t *Table) HalfPageUp() {
	t.SelectedRow = t.SelectedRow - (t.Y-2)/2
	t.calcPos()
}

func (t *Table) HalfPageDown() {
	t.SelectedRow = t.SelectedRow + (t.Y-2)/2
	t.calcPos()
}

func (t *Table) PageUp() {
	t.SelectedRow -= (t.Y - 2)
	t.calcPos()
}

func (t *Table) PageDown() {
	t.SelectedRow += (t.Y - 2)
	t.calcPos()
}

func (t *Table) Click(x, y int) {
	x = x - t.XOffset
	y = y - t.YOffset
	if (x > 0 && x <= t.X) && (y > 0 && y <= t.Y) {
		t.SelectedRow = (t.TopRow + y) - 2
		t.calcPos()
	}
}
