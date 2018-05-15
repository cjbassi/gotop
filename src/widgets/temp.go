package widgets

// Temp is too customized to inherit from a generic widget so we create a customized one here.
// Temp defines its own Buffer method directly.

import (
	"fmt"
	"sort"
	"time"

	ui "github.com/cjbassi/termui"
)

type Temp struct {
	*ui.Block
	interval  time.Duration
	Data      map[string]int
	Threshold int
	TempLow   ui.Color
	TempHigh  ui.Color
}

func NewTemp() *Temp {
	self := &Temp{
		Block:     ui.NewBlock(),
		interval:  time.Second * 5,
		Data:      make(map[string]int),
		Threshold: 80, // temp at which color should change
	}
	self.Label = "Temperatures"

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

// Buffer implements ui.Bufferer interface and renders the widget.
func (self *Temp) Buffer() *ui.Buffer {
	buf := self.Block.Buffer()

	var keys []string
	for key := range self.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for y, key := range keys {
		if y+1 > self.Y {
			break
		}

		fg := self.TempLow
		if self.Data[key] >= self.Threshold {
			fg = self.TempHigh
		}

		s := ui.MaxString(key, (self.X - 4))
		buf.SetString(1, y+1, s, self.Fg, self.Bg)
		buf.SetString(self.X-2, y+1, fmt.Sprintf("%2dC", self.Data[key]), fg, self.Bg)
	}

	return buf
}
