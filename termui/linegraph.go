package termui

import (
	"image"
	"sort"
	"strconv"
	"unicode"

	. "github.com/gizak/termui/v3"
	drawille "github.com/xxxserxxx/gotop/v4/termui/drawille-go"
)

// LineGraph draws a graph like this ⣀⡠⠤⠔⣁ of data points.
type LineGraph struct {
	*Block

	// Data is a size-managed data set for the graph. Each entry is a line;
	// each sub-array are points in the line. The maximum size of the
	// sub-arrays is controlled by the size of the canvas. This
	// array is **not** thread-safe. Do not modify this array, or it's
	// sub-arrays in threads different than the thread that calls `Draw()`
	Data map[string][]float64
	// The labels drawn on the graph for each of the lines; the key is shared
	// by Data; the value is the text that will be rendered.
	Labels map[string]string

	HorizontalScale int

	LineColors       map[string]Color
	LabelStyles      map[string]Modifier
	DefaultLineColor Color

	seriesList numbered
}

func NewLineGraph() *LineGraph {
	return &LineGraph{
		Block: NewBlock(),

		Data:   make(map[string][]float64),
		Labels: make(map[string]string),

		HorizontalScale: 5,

		LineColors:  make(map[string]Color),
		LabelStyles: make(map[string]Modifier),
	}
}

func (self *LineGraph) Draw(buf *Buffer) {
	self.Block.Draw(buf)
	// we render each data point on to the canvas then copy over the braille to the buffer at the end
	// fyi braille characters have 2x4 dots for each character
	c := drawille.NewCanvas()
	// used to keep track of the braille colors until the end when we render the braille to the buffer
	colors := make([][]Color, self.Inner.Dx()+2)
	for i := range colors {
		colors[i] = make([]Color, self.Inner.Dy()+2)
	}

	if len(self.seriesList) != len(self.Data) {
		// sort the series so that overlapping data will overlap the same way each time
		self.seriesList = make(numbered, len(self.Data))
		i := 0
		for seriesName := range self.Data {
			self.seriesList[i] = seriesName
			i++
		}
		sort.Sort(self.seriesList)
	}

	// draw lines in reverse order so that the first color defined in the colorscheme is on top
	for i := len(self.seriesList) - 1; i >= 0; i-- {
		seriesName := self.seriesList[i]
		seriesData := self.Data[seriesName]
		seriesLineColor, ok := self.LineColors[seriesName]
		if !ok {
			seriesLineColor = self.DefaultLineColor
			self.LineColors[seriesName] = seriesLineColor
		}

		// coordinates of last point
		lastY, lastX := -1, -1
		// assign colors to `colors` and lines/points to the canvas
		dx := self.Inner.Dx()
		for i := len(seriesData) - 1; i >= 0; i-- {
			x := ((dx + 1) * 2) - 1 - (((len(seriesData) - 1) - i) * self.HorizontalScale)
			y := ((self.Inner.Dy() + 1) * 4) - 1 - int((float64((self.Inner.Dy())*4)-1)*(seriesData[i]/100))
			if x < 0 {
				// render the line to the last point up to the wall
				if x > -self.HorizontalScale {
					for _, p := range drawille.Line(lastX, lastY, x, y) {
						if p.X > 0 {
							c.Set(p.X, p.Y)
							colors[p.X/2][p.Y/4] = seriesLineColor
						}
					}
				}
				if len(seriesData) > 4*dx {
					self.Data[seriesName] = seriesData[dx-1:]
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
					buf.SetCell(
						NewCell(char, NewStyle(colors[x][y])),
						image.Pt(self.Inner.Min.X+x-1, self.Inner.Min.Y+y-1),
					)
				}
			}
		}
	}

	// renders key/label ontop
	maxWid := 0
	xoff := 0 // X offset for additional columns of text
	yoff := 0 // Y offset for resetting column to top of widget
	for i, seriesName := range self.seriesList {
		if yoff+i+2 > self.Inner.Dy() {
			xoff += maxWid + 2
			yoff = -i
			maxWid = 0
		}
		seriesLineColor, ok := self.LineColors[seriesName]
		if !ok {
			seriesLineColor = self.DefaultLineColor
		}
		seriesLabelStyle, ok := self.LabelStyles[seriesName]
		if !ok {
			seriesLabelStyle = ModifierClear
		}

		// render key ontop, but let braille be drawn over space characters
		str := seriesName + " " + self.Labels[seriesName]
		if len(str) > maxWid {
			maxWid = len(str)
		}
		for k, char := range str {
			if char != ' ' {
				buf.SetCell(
					NewCell(char, NewStyle(seriesLineColor, ColorClear, seriesLabelStyle)),
					image.Pt(xoff+self.Inner.Min.X+2+k, yoff+self.Inner.Min.Y+i+1),
				)
			}
		}

	}
}

// A string containing an integer
type numbered []string

func (n numbered) Len() int      { return len(n) }
func (n numbered) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func (n numbered) Less(i, j int) bool {
	ars := n[i]
	brs := n[j]
	// iterate over the runes in string A
	var ai int
	for ai = 0; ai < len(ars); ai++ {
		ar := ars[ai]
		// if brs has been equal to ars so far, but is shorter, than brs is less
		if ai >= len(brs) {
			return false
		}
		br := brs[ai]
		// handle numbers if we find them
		if unicode.IsDigit(rune(ar)) {
			// if ar is a digit and br isn't, then ars is less
			if !unicode.IsDigit(rune(brs[ai])) {
				return true
			}
			// now get the range for the full numbers, which is the sequence of consecutive digits from each
			ae := ai + 1
			be := ae
			for ; ae < len(ars) && unicode.IsDigit(rune(ars[ae])); ae++ {
			}
			for ; be < len(brs) && unicode.IsDigit(rune(brs[be])); be++ {
			}
			if be < ae {
				return false
			}
			adigs, err := strconv.Atoi(string(ars[ai:ae]))
			if err != nil {
				return true
			}
			bdigs, err := strconv.Atoi(string(brs[ai:be]))
			if err != nil {
				return true
			}
			// The numbers are equal; continue with the rest of the string
			if adigs == bdigs {
				ai = ae - 1
				continue
			}
			return adigs < bdigs
		} else if unicode.IsDigit(rune(br)) {
			return false
		}
		if ar == br {
			continue
		}
		// No digits, so compare letters
		if ar < br {
			return true
		}
		return false
	}
	if ai <= len(brs) {
		return true
	}
	return false
}
