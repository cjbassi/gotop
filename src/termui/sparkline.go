package termui

import (
	"image"
	"log"

	. "github.com/gizak/termui"
)

// Sparkline is like: ▅▆▂▂▅▇▂▂▃▆▆▆▅▃. The data points should be non-negative integers.
type Sparkline struct {
	Data       []int
	Title1     string
	Title2     string
	TitleColor Color
	LineColor  Color
}

// Sparklines is a renderable widget which groups together the given sparklines.
type Sparklines struct {
	*Block
	Lines []*Sparkline
}

// Add appends a given Sparkline to the *Sparklines.
func (self *Sparklines) Add(sl Sparkline) {
	self.Lines = append(self.Lines, &sl)
}

// NewSparkline returns an unrenderable single sparkline that intended to be added into a Sparklines.
func NewSparkline() *Sparkline {
	return &Sparkline{}
}

// NewSparklines return a new *Sparklines with given Sparklines, you can always add a new Sparkline later.
func NewSparklines(ss ...*Sparkline) *Sparklines {
	return &Sparklines{
		Block: NewBlock(),
		Lines: ss,
	}
}

func (self *Sparklines) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	lc := len(self.Lines) // lineCount

	// renders each sparkline and its titles
	for i, line := range self.Lines {

		// prints titles
		title1Y := self.Inner.Min.Y + 1 + (self.Inner.Dy()/lc)*i
		title2Y := self.Inner.Min.Y + 2 + (self.Inner.Dy()/lc)*i
		title1 := TrimString(line.Title1, self.Inner.Dx())
		title2 := TrimString(line.Title2, self.Inner.Dx())
		if self.Inner.Dy() > 5 {
			buf.SetString(
				title1,
				NewStyle(line.TitleColor, ColorClear, ModifierBold),
				image.Pt(self.Inner.Min.X, title1Y),
			)
		}
		if self.Inner.Dy() > 6 {
			buf.SetString(
				title2,
				NewStyle(line.TitleColor, ColorClear, ModifierBold),
				image.Pt(self.Inner.Min.X, title2Y),
			)
		}

		sparkY := (self.Inner.Dy() / lc) * (i + 1)
		// finds max data in current view used for relative heights
		max := 1
		for i := len(line.Data) - 1; i >= 0 && self.Inner.Dx()-((len(line.Data)-1)-i) >= 1; i-- {
			if line.Data[i] > max {
				max = line.Data[i]
			}
		}
		// prints sparkline
		for x := self.Inner.Dx(); x >= 1; x-- {
			char := BARS[0]
			if (self.Inner.Dx() - x) < len(line.Data) {
				offset := self.Inner.Dx() - x
				curItem := line.Data[(len(line.Data)-1)-offset]
				percent := float64(curItem) / float64(max)
				index := int(percent * 7)
				if index < 0 || index >= len(BARS) {
					log.Printf(
						"invalid sparkline data value. index: %v, percent: %v, curItem: %v, offset: %v",
						index, percent, curItem, offset,
					)
				} else {
					char = BARS[index]
				}
			}
			buf.SetCell(
				NewCell(char, NewStyle(line.LineColor)),
				image.Pt(self.Inner.Min.X+x-1, self.Inner.Min.Y+sparkY-1),
			)
		}
	}
}
