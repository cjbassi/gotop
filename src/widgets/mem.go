package widgets

import (
	"time"

	ui "github.com/cjbassi/gotop/src/termui"
)

type MemWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
}

func NewMemWidget(updateInterval time.Duration, horizontalScale int) *MemWidget {
	self := &MemWidget{
		LineGraph:      ui.NewLineGraph(),
		updateInterval: updateInterval,
	}
	self.Title = " Memory Usage "
	self.HorizontalScale = horizontalScale
	self.Data["Main"] = []float64{0}
	self.Data["Swap"] = []float64{0}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}
