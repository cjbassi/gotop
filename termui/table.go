package termui

import (
	"os/exec"
	"strings"
)

// Table tracks all the attributes of a Table instance
type Table struct {
	*Block
	Header    []string
	Rows      [][]string
	Fg        Color
	Bg        Color
	Cursor    Color
	UniqueCol int
	pid       string
	selected  int
	topRow    int
}

// NewTable returns a new Table instance
func NewTable() *Table {
	return &Table{
		Block:     NewBlock(),
		Fg:        Theme.Fg,
		Bg:        Theme.Bg,
		Cursor:    Theme.TableCursor,
		selected:  0,
		topRow:    0,
		UniqueCol: 0,
	}
}

// Buffer ...
func (t *Table) Buffer() *Buffer {
	buf := t.Block.Buffer()

	if t.topRow > len(t.Rows)-(t.Y-1) {
		t.topRow = len(t.Rows) - (t.Y - 1)
	}

	// calculate gap size based on total width
	gap := 3
	if t.X < 50 {
		gap = 1
	} else if t.X < 75 {
		gap = 2
	}

	cw := []int{5, 10, 4, 4} // cellWidth
	cp := []int{             // cellPosition
		gap,
		gap + cw[0] + gap,
		t.X - gap - cw[3] - gap - cw[2],
		t.X - gap - cw[3],
	}

	// total width requires by all 4 columns
	contentWidth := gap + cw[0] + gap + cw[1] + gap + cw[2] + gap + cw[3] + gap
	render := 4 // number of columns to iterate through

	// removes CPU and MEM if there isn't enough room
	if t.X < (contentWidth - gap - cw[3]) {
		render = 2
	} else if t.X < contentWidth {
		cp[2] = cp[3]
		render = 3
	}

	// print header
	for i := 0; i < render; i++ {
		r := MaxString(t.Header[i], t.X-6)
		buf.SetString(cp[i], 1, r, t.Fg|AttrBold, t.Bg)
	}

	// prints each row
	// for y, row := range t.Rows {
	// for y := t.topRow; y <= t.topRow+t.Y; y++ {
	for rowNum := t.topRow; rowNum < t.topRow+t.Y-1 && rowNum < len(t.Rows); rowNum++ {
		row := t.Rows[rowNum]
		y := (rowNum + 2) - t.topRow

		// cursor
		bg := t.Bg
		if (t.pid == "" && rowNum == t.selected) || (t.pid != "" && t.pid == row[t.UniqueCol]) {
			bg = t.Cursor
			for i := 0; i < render; i++ {
				buf.SetString(1, y, strings.Repeat(" ", t.X), t.Fg, bg)
			}
			t.pid = row[t.UniqueCol]
			t.selected = rowNum
		}

		// prints each string
		for i := 0; i < render; i++ {
			r := MaxString(row[i], t.X-6)
			buf.SetString(cp[i], y, r, t.Fg, bg)
		}
	}

	return buf
}

////////////////////////////////////////////////////////////////////////////////

func (t *Table) calcPos() {
	t.pid = ""

	if t.selected < 0 {
		t.selected = 0
	}
	if t.selected < t.topRow {
		t.topRow = t.selected
	}

	if t.selected > len(t.Rows)-1 {
		t.selected = len(t.Rows) - 1
	}
	if t.selected > t.topRow+(t.Y-2) {
		t.topRow = t.selected - (t.Y - 2)
	}
}

func (t *Table) Up() {
	t.selected -= 1
	t.calcPos()
}

func (t *Table) Down() {
	t.selected += 1
	t.calcPos()
}

func (t *Table) Top() {
	t.selected = 0
	t.calcPos()
}

func (t *Table) Bottom() {
	t.selected = len(t.Rows) - 1
	t.calcPos()
}

func (t *Table) HalfPageUp() {
	t.selected = t.selected - (t.Y-2)/2
	t.calcPos()
}

func (t *Table) HalfPageDown() {
	t.selected = t.selected + (t.Y-2)/2
	t.calcPos()
}

func (t *Table) PageUp() {
	t.selected -= (t.Y - 2)
	t.calcPos()
}

func (t *Table) PageDown() {
	t.selected += (t.Y - 2)
	t.calcPos()
}

func (t *Table) Click(x, y int) {
	x = x - t.XOffset
	y = y - t.YOffset
	if (x > 0 && x <= t.X) && (y > 0 && y <= t.Y) {
		t.selected = (t.topRow + y) - 2
		t.calcPos()
	}
}

func (t *Table) Kill() {
	t.pid = ""
	command := "kill"
	if t.UniqueCol == 1 {
		command = "pkill"
	}
	cmd := exec.Command(command, t.Rows[t.selected][t.UniqueCol])
	cmd.Start()
}
