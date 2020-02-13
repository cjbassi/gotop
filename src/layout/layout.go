// Copyright 2020 The gotop Authors Licensed under terms of the LICENSE file in this repository.
package layout

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/config"
	"github.com/cjbassi/gotop/src/widgets"
	ui "github.com/gizak/termui/v3"
)

type layout struct {
	Rows [][]widgetRule
}

type widgetRule struct {
	Widget string
	Weight float64
	Height int
}

type MyGrid struct {
	*ui.Grid
	Lines []widgets.Scalable
	Proc  *widgets.ProcWidget
}

var widgetNames []string = []string{"cpu", "disk", "mem", "temp", "net", "procs", "batt"}

// The syntax for the layout specification is:
// ```
// (rowspan:)?widget(/weight)?
// ```
// 1. Each line is a row
// 2. Empty lines are skipped
// 3. Spaces are compressed
// 4. Legal widget names are: cpu, disk, mem, temp, batt, net, procs
// 5. Names are not case sensitive
// 4. The simplest row is a single widget, by name, e.g.
//    ```
//    cpu
//    ```
// 5. Widgets with no weights have a weight of 1.
// 6. If multiple widgets are put on a row with no weights, they will all have
//    the same width.
// 7. Weights are integers
// 8. A widget will have a width proportional to its weight divided by the
//    total weight count of the row. E.g.,
//    ```
//    cpu      net
//    disk/2   mem/4
//    ```
//    The first row will have two widgets: the CPU and network widgets; each
//    will be 50% of the total width wide.  The second row will have two
//    widgets: disk and memory; the first will be 2/6 ~= 33% wide, and the
//    second will be 5/7 ~= 67% wide (or, memory will be twice as wide as disk).
// 9. If prefixed by a number and colon, the widget will span that number of
//    rows downward. E.g.
//    ```
//    2:cpu
//    mem
//    ```
//    The CPU widget will be twice as high as the memory widget.  Similarly,
//    ```
//    mem   2:cpu
//    net
//    ```
//    memory and network will be in the same row as CPU, one over the other,
//    and each half as high as CPU.
// 10. Negative, 0, or non-integer weights will be recorded as "1".  Same for row spans.
// 11. Unrecognized widgets will cause the application to abort.
// 12. In rows with multi-row spanning widgets **and** weights, weights in
//     lower rows are ignored.  Put the weight on the widgets in that row, not
//     in later (spanned) rows.
func ParseLayout(i io.Reader) layout {
	r := bufio.NewScanner(i)
	rv := layout{Rows: make([][]widgetRule, 0)}
	var lineNo int
	for r.Scan() {
		l := strings.TrimSpace(r.Text())
		if l == "" {
			continue
		}
		row := make([]widgetRule, 0)
		ws := strings.Fields(l)
		weightTotal := 0
		for _, w := range ws {
			wr := widgetRule{Weight: 1}
			ks := strings.Split(w, "/")
			rs := strings.Split(ks[0], ":")
			var wid string
			if len(rs) > 1 {
				v, e := strconv.Atoi(rs[0])
				if e != nil {
					log.Printf("Layout error on line %d: format must be INT:STRING/INT. Error parsing %s as a int. Word was %s. Using a row height of 1.", lineNo, rs[0], w)
					v = 1
				}
				if v < 1 {
					v = 1
				}
				wr.Height = v
				wid = rs[1]
			} else {
				wr.Height = 1
				wid = rs[0]
			}
			wr.Widget = strings.ToLower(wid)
			if len(ks) > 1 {
				weight, e := strconv.Atoi(ks[1])
				if e != nil {
					log.Printf("Layout error on line %d: format must be STRING/INT. Error parsing %s as a int. Word was %s. Using a weight of 1 for widget.", lineNo, ks[1], w)
					weight = 1
				}
				if weight < 1 {
					weight = 1
				}
				wr.Weight = float64(weight)
				if len(ks) > 2 {
					log.Printf("Layout warning on line %d: too many '/' in word %s; ignoring extra junk.", lineNo, w)
				}
				weightTotal += weight
			} else {
				weightTotal += 1
			}
			row = append(row, wr)
		}
		// Prevent tricksy users from breaking their own computers
		if weightTotal <= 1 {
			weightTotal = 1
		}
		for i, w := range row {
			row[i].Weight = w.Weight / float64(weightTotal)
		}
		rv.Rows = append(rv.Rows, row)
	}
	return rv
}

func Layout(wl layout, c config.Config) (*MyGrid, error) {
	log.Printf("laying out %v", wl)
	var rows [][]interface{}
	var lines []widgets.Scalable
	var mouser *widgets.ProcWidget
	for _, rowDef := range wl.Rows {
		var gi []interface{}
		var w interface{}
		for _, widRule := range rowDef {
			switch widRule.Widget {
			case "cpu":
				cpu := widgets.NewCpuWidget(c.UpdateInterval, c.GraphHorizontalScale, c.AverageLoad, c.PercpuLoad)
				var keys []string
				for key := range cpu.Data {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				i := 0
				for _, v := range keys {
					if i >= len(c.Colorscheme.CPULines) {
						// assuming colorscheme for CPU lines is not empty
						i = 0
					}
					color := c.Colorscheme.CPULines[i]
					cpu.LineColors[v] = ui.Color(color)
					i++
				}
				lines = append(lines, cpu)
				w = cpu
			case "disk":
				w = widgets.NewDiskWidget()
			case "mem":
				m := widgets.NewMemWidget(c.UpdateInterval, c.GraphHorizontalScale)
				m.LineColors["Main"] = ui.Color(c.Colorscheme.MainMem)
				m.LineColors["Swap"] = ui.Color(c.Colorscheme.SwapMem)
				lines = append(lines, m)
				w = m
			case "temp":
				t := widgets.NewTempWidget(c.TempScale)
				t.TempLowColor = ui.Color(c.Colorscheme.TempLow)
				t.TempHighColor = ui.Color(c.Colorscheme.TempHigh)
				w = t
			case "net":
				n := widgets.NewNetWidget(c.NetInterface)
				n.Lines[0].LineColor = ui.Color(c.Colorscheme.Sparkline)
				n.Lines[0].TitleColor = ui.Color(c.Colorscheme.BorderLabel)
				n.Lines[1].LineColor = ui.Color(c.Colorscheme.Sparkline)
				n.Lines[1].TitleColor = ui.Color(c.Colorscheme.BorderLabel)
				w = n
			case "procs":
				p := widgets.NewProcWidget()
				p.CursorColor = ui.Color(c.Colorscheme.ProcCursor)
				mouser = p
				w = p
			case "batt":
				b := widgets.NewBatteryWidget(c.GraphHorizontalScale)
				var battKeys []string
				for key := range b.Data {
					battKeys = append(battKeys, key)
				}
				sort.Strings(battKeys)
				i := 0 // Re-using variable from CPU
				for _, v := range battKeys {
					if i >= len(c.Colorscheme.BattLines) {
						// assuming colorscheme for battery lines is not empty
						i = 0
					}
					color := c.Colorscheme.BattLines[i]
					b.LineColors[v] = ui.Color(color)
					i++
				}
				w = b
			default:
				return nil, errors.New(fmt.Sprintf("Invalid widget name %s.  Must be one of %v", widRule.Widget, widgetNames))
			}
			gi = append(gi, ui.NewCol(widRule.Weight, w))
		}
		if len(gi) > 0 {
			rows = append(rows, gi)
		} else {
			log.Printf("WARN: no rows created from %v", rowDef)
		}
	}
	var rgs []interface{}
	rowHeight := 1.0 / float64(len(rows))
	for _, r := range rows {
		rgs = append(rgs, ui.NewRow(rowHeight, r...))
	}
	grid := &MyGrid{ui.NewGrid(), make([]widgets.Scalable, 0), mouser}
	grid.Set(rgs...)
	grid.Lines = lines
	return grid, nil
}
