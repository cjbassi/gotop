package widgets

// Temp is too customized to inherit from a generic widget so we create a customized one here.
// Temp defines its own Buffer method directly.

import (
	"fmt"
	"image"
	"sort"
	"time"

	ui "github.com/gizak/termui"
)

type Temp struct {
	*ui.Block
	interval   time.Duration
	Data       map[string]int
	Threshold  int
	TempLow    ui.Attribute
	TempHigh   ui.Attribute
	Fahrenheit bool
}

func NewTemp(fahrenheit bool) *Temp {
	self := &Temp{
		Block:     ui.NewBlock(),
		interval:  time.Second * 5,
		Data:      make(map[string]int),
		Threshold: 80, // temp at which color should change
	}
	self.Title = " Temperatures "

	if fahrenheit {
		self.Fahrenheit = true
		self.Threshold = int(self.Threshold*9/5 + 32)
	}

	self.update()

	go func() {
		ticker := time.NewTicker(self.interval)
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Temp) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	var keys []string
	for key := range self.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for y, key := range keys {
		if y+1 > self.Inner.Dy() {
			break
		}

		fg := self.TempLow
		if self.Data[key] >= self.Threshold {
			fg = self.TempHigh
		}

		s := ui.TrimString(key, (self.Inner.Dx() - 4))
		buf.SetString(s,
			image.Pt(self.Inner.Min.X, self.Inner.Min.Y+y),
			ui.Theme.Default,
		)
		if self.Fahrenheit {
			buf.SetString(
				fmt.Sprintf("%3dF", self.Data[key]),
				image.Pt(self.Inner.Dx()-3, y+1),
				ui.AttrPair{fg, -1},
			)
		} else {
			buf.SetString(
				fmt.Sprintf("%3dC", self.Data[key]),
				image.Pt(self.Inner.Max.X-4, self.Inner.Min.Y+y),
				ui.AttrPair{fg, -1},
			)
		}
	}
}
