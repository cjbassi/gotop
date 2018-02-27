package termui

import (
	"fmt"
	"sort"

	drawille "github.com/cjbassi/drawille-go"
)

// LineGraph implements a graph of data points.
type LineGraph struct {
	*Block
	Data      map[string][]float64
	LineColor map[string]Color

	DefaultLineColor Color
}

// NewLineGraph returns a new LineGraph with current theme.
func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block:     NewBlock(),
		Data:      make(map[string][]float64),
		LineColor: make(map[string]Color),

		DefaultLineColor: Theme.LineGraph,
	}
}

// renderPoints plots and interpolates data points.
func (lc *LineGraph) Buffer() *Buffer {
	buf := lc.Block.Buffer()
	c := drawille.NewCanvas()
	// used to keep track of colors but not write them to the buffer until the end
	colors := make([][]Color, lc.X+2)
	for i := range colors {
		colors[i] = make([]Color, lc.Y+2)
	}

	// Sort the series so that overlapping data will overlap the same way each time
	seriesList := make([]string, len(lc.Data))
	i := 0
	for seriesName := range lc.Data {
		seriesList[i] = seriesName
		i++
	}
	sort.Strings(seriesList)

	// draw lines in reverse order so the first one is on top
	for i := len(seriesList) - 1; i >= 0; i-- {
		seriesName := seriesList[i]
		seriesData := lc.Data[seriesName]
		seriesLineColor, ok := lc.LineColor[seriesName]
		if !ok {
			seriesLineColor = lc.DefaultLineColor
		}

		// coordinates of last point
		lastY, lastX := -1, -1
		// assign colors to `colors` and lines/points to the canvas
		for i := len(seriesData) - 1; i >= 0; i-- {
			x := ((lc.X + 1) * 2) - 1 - (((len(seriesData) - 1) - i) * 5)
			y := ((lc.Y + 1) * 4) - 1 - int((float64((lc.Y)*4)-1)*(seriesData[i]/100))
			// stop rendering at the left-most wall
			if x < 0 {
				break
			}
			if lastY == -1 { // if this is the first point
				c.Set(x, y)
				colors[x/2][y/4] = seriesLineColor // divide by 2 and 4 due to 2x4 dots in braille characters
			} else {
				c.DrawLine(lastX, lastY, x, y)
				for _, p := range drawille.Line(lastX, lastY, x, y) {
					colors[p.X/2][p.Y/4] = seriesLineColor
				}
			}
			lastX, lastY = x, y
		}

		// copy braille and colors to buffer
		for y, line := range c.Rows(c.MinX(), c.MinY(), c.MaxX(), c.MaxY()) {
			for x, char := range line {
				x /= 3 // idk why but it works
				if x == 0 {
					continue
				}
				if char != 10240 { // empty braille character
					buf.SetCell(x, y, Cell{char, colors[x][y], lc.Bg})
				}
			}
		}
	}

	for j, seriesName := range seriesList {
		seriesData := lc.Data[seriesName]
		seriesLineColor, ok := lc.LineColor[seriesName]
		if !ok {
			seriesLineColor = lc.DefaultLineColor
		}

		// Render key ontop, but let braille be drawn between words
		str := fmt.Sprintf("%s %3.0f%%", seriesName, seriesData[len(seriesData)-1])
		for k, char := range str {
			if char != ' ' {
				buf.SetCell(3+k, j+2, Cell{char, seriesLineColor, lc.Bg})
			}
		}

	}

	return buf
}
