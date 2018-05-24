package termui

import (
	"sort"

	drawille "github.com/cjbassi/drawille-go"
)

// LineGraph implements a line graph of data points.
type LineGraph struct {
	*Block
	Data      map[string][]float64
	LineColor map[string]Color
	Zoom      int
	Labels    map[string]string

	DefaultLineColor Color
}

// NewLineGraph returns a new LineGraph with current theme.
func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block:     NewBlock(),
		Data:      make(map[string][]float64),
		LineColor: make(map[string]Color),
		Labels:    make(map[string]string),
		Zoom:      5,

		DefaultLineColor: Theme.LineGraph,
	}
}

// Buffer implements Bufferer interface.
func (self *LineGraph) Buffer() *Buffer {
	buf := self.Block.Buffer()
	// we render each data point on to the canvas then copy over the braille to the buffer at the end
	// fyi braille characters have 2x4 dots for each character
	c := drawille.NewCanvas()
	// used to keep track of the braille colors until the end when we render the braille to the buffer
	colors := make([][]Color, self.X+2)
	for i := range colors {
		colors[i] = make([]Color, self.Y+2)
	}

	// sort the series so that overlapping data will overlap the same way each time
	seriesList := make([]string, len(self.Data))
	i := 0
	for seriesName := range self.Data {
		seriesList[i] = seriesName
		i++
	}
	sort.Strings(seriesList)

	// draw lines in reverse order so that the first color defined in the colorscheme is on top
	for i := len(seriesList) - 1; i >= 0; i-- {
		seriesName := seriesList[i]
		seriesData := self.Data[seriesName]
		seriesLineColor, ok := self.LineColor[seriesName]
		if !ok {
			seriesLineColor = self.DefaultLineColor
		}

		// coordinates of last point
		lastY, lastX := -1, -1
		// assign colors to `colors` and lines/points to the canvas
		for i := len(seriesData) - 1; i >= 0; i-- {
			x := ((self.X + 1) * 2) - 1 - (((len(seriesData) - 1) - i) * self.Zoom)
			y := ((self.Y + 1) * 4) - 1 - int((float64((self.Y)*4)-1)*(seriesData[i]/100))
			if x < 0 {
				// render the line to the last point up to the wall
				if x > 0-self.Zoom {
					for _, p := range drawille.Line(lastX, lastY, x, y) {
						if p.X > 0 {
							c.Set(p.X, p.Y)
							colors[p.X/2][p.Y/4] = seriesLineColor
						}
					}
				}
				break
			}
			if lastY == -1 { // if this is the first point
				c.Set(x, y)
				colors[x/2][y/4] = seriesLineColor
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
					buf.SetCell(x, y, Cell{char, colors[x][y], self.Bg})
				}
			}
		}
	}

	// renders key/label ontop
	for j, seriesName := range seriesList {
		seriesLineColor, ok := self.LineColor[seriesName]
		if !ok {
			seriesLineColor = self.DefaultLineColor
		}

		// render key ontop, but let braille be drawn over space characters
		str := seriesName + " " + self.Labels[seriesName]
		for k, char := range str {
			if char != ' ' {
				buf.SetCell(3+k, j+2, Cell{char, seriesLineColor, self.Bg})
			}
		}

	}

	return buf
}
