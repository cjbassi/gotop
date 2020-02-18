// Copyright 2020 The gotop Authors Licensed under terms of the LICENSE file in this repository.
package layout

import (
	"log"
	"sort"

	"github.com/xxxserxxx/gotop"
	"github.com/xxxserxxx/gotop/widgets"

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

func Layout(wl layout, c gotop.Config) (*MyGrid, error) {
	rowDefs := wl.Rows
	uiRows := make([]ui.GridItem, 0)
	numRows := countNumRows(wl.Rows)
	var uiRow ui.GridItem
	for len(rowDefs) > 0 {
		uiRow, rowDefs = processRow(c, numRows, rowDefs)
		uiRows = append(uiRows, uiRow)
	}
	rgs := make([]interface{}, 0)
	for _, ur := range uiRows {
		ur.HeightRatio = ur.HeightRatio / float64(numRows)
		rgs = append(rgs, ur)
	}
	grid := &MyGrid{ui.NewGrid(), nil, nil}
	grid.Set(rgs...)
	grid.Lines = deepFindScalable(rgs)
	grid.Proc = deepFindProc(uiRows)
	return grid, nil
}

// processRow eats a single row from the input list of rows and returns a UI
// row (GridItem) representation of the specification, along with a slice
// without that row.
//
// It does more than that, actually, because it may consume more than one row
// if there's a row span widget in the row; in this case, it'll consume as many
// rows as the largest row span object in the row, and produce an uber-row
// containing all that stuff. It returns a slice without the consumed elements.
func processRow(c gotop.Config, numRows int, rowDefs [][]widgetRule) (ui.GridItem, [][]widgetRule) {
	// Recursive function #3.  See the comment in deepFindProc.
	if len(rowDefs) < 1 {
		return ui.GridItem{}, [][]widgetRule{}
	}
	// The height of the tallest widget in this row; the number of rows that
	// will be consumed, and the overall height of the row that will be
	// produced.
	maxHeight := countMaxHeight([][]widgetRule{rowDefs[0]})
	var processing [][]widgetRule
	if maxHeight < len(rowDefs) {
		processing = rowDefs[0:maxHeight]
		rowDefs = rowDefs[maxHeight:]
	} else {
		processing = rowDefs[0:]
		rowDefs = [][]widgetRule{}
	}
	var colWeights []float64
	var columns [][]interface{}
	numCols := len(processing[0])
	if numCols < 1 {
		numCols = 1
	}
	for _, rd := range processing[0] {
		colWeights = append(colWeights, rd.Weight)
		columns = append(columns, make([]interface{}, 0))
	}
	colHeights := make([]int, numCols)
	for _, rds := range processing {
		for _, rd := range rds {
			for j, ch := range colHeights {
				if ch+rd.Height <= maxHeight {
					widget := makeWidget(c, rd)
					columns[j] = append(columns[j], ui.NewRow(float64(rd.Height)/float64(maxHeight), widget))
					colHeights[j] += rd.Height
					break
				}
			}
		}
	}
	var uiColumns []interface{}
	for i, widgets := range columns {
		if len(widgets) > 0 {
			uiColumns = append(uiColumns, ui.NewCol(float64(colWeights[i]), widgets...))
		}
	}

	return ui.NewRow(1.0/float64(numRows), uiColumns...), rowDefs
}

type Metric interface {
	EnableMetric()
}

func makeWidget(c gotop.Config, widRule widgetRule) interface{} {
	var w Metric
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
		w = cpu
	case "disk":
		dw := widgets.NewDiskWidget()
		w = dw
	case "mem":
		m := widgets.NewMemWidget(c.UpdateInterval, c.GraphHorizontalScale)
		m.LineColors["Main"] = ui.Color(c.Colorscheme.MainMem)
		m.LineColors["Swap"] = ui.Color(c.Colorscheme.SwapMem)
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
		log.Printf("Invalid widget name %s.  Must be one of %v", widRule.Widget, widgetNames)
		return ui.NewBlock()
	}
	if c.ExportPort != "" {
		w.EnableMetric()
	}
	return w
}

func countNumRows(rs [][]widgetRule) int {
	var ttl int
	for len(rs) > 0 {
		ttl += 1
		line := rs[0]
		h := 1
		for _, c := range line {
			if c.Height > h {
				h = c.Height
			}
		}
		if h < len(rs) {
			rs = rs[h:]
		} else {
			break
		}
	}
	return ttl
}

// Counts the height of the window so rows can be proportionally scaled.
func countMaxHeight(rs [][]widgetRule) int {
	var ttl int
	for len(rs) > 0 {
		line := rs[0]
		h := 1
		for _, c := range line {
			if c.Height > h {
				h = c.Height
			}
		}
		ttl += h
		if h < len(rs) {
			rs = rs[h:]
		} else {
			break
		}
	}
	return ttl
}

// deepFindProc looks in the UI widget tree for the ProcWidget,
// and returns it if found or nil if not.
func deepFindProc(gs interface{}) *widgets.ProcWidget {
	// Recursive function #1.  Recursion is OK here because the number
	// of UI elements, even in a very complex UI, is going to be
	// relatively small.
	t, ok := gs.(ui.GridItem)
	if ok {
		return deepFindProc(t.Entry)
	}
	es, ok := gs.([]ui.GridItem)
	if ok {
		for _, g := range es {
			v := deepFindProc(g)
			if v != nil {
				return v
			}
		}
	}
	fs, ok := gs.([]interface{})
	if ok {
		for _, g := range fs {
			v := deepFindProc(g)
			if v != nil {
				return v
			}
		}
	}
	p, ok := gs.(*widgets.ProcWidget)
	if ok {
		return p
	}
	return nil
}

// deepFindScalable looks in the UI widget tree for Scalable widgets,
// and returns them if found or an empty slice if not.
func deepFindScalable(gs interface{}) []widgets.Scalable {
	// Recursive function #1.  See the comment in deepFindProc.
	t, ok := gs.(ui.GridItem)
	if ok {
		return deepFindScalable(t.Entry)
	}
	es, ok := gs.([]ui.GridItem)
	rvs := make([]widgets.Scalable, 0)
	if ok {
		for _, g := range es {
			vs := deepFindScalable(g)
			rvs = append(rvs, vs...)
		}
		return rvs
	}
	fs, ok := gs.([]interface{})
	if ok {
		for _, g := range fs {
			vs := deepFindScalable(g)
			rvs = append(rvs, vs...)
		}
		return rvs
	}
	p, ok := gs.(widgets.Scalable)
	if ok {
		rvs = append(rvs, p)
	}
	return rvs
}
