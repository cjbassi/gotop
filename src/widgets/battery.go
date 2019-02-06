package widgets

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/cjbassi/battery"

	ui "github.com/cjbassi/gotop/src/termui"
)

type Batt struct {
	*ui.LineGraph
	interval time.Duration
}

func NewBatt(renderLock *sync.RWMutex, horizontalScale int) *Batt {
	self := &Batt{
		LineGraph: ui.NewLineGraph(),
		interval:  time.Minute,
	}
	self.Title = " Battery Status "
	self.HorizontalScale = horizontalScale

	// intentional duplicate
	self.update()
	self.update()

	go func() {
		for range time.NewTicker(self.interval).C {
			renderLock.RLock()
			self.update()
			renderLock.RUnlock()
		}
	}()

	return self
}

func mkId(i int) string {
	return "Batt" + strconv.Itoa(i)
}

func (self *Batt) update() {
	batts, err := battery.GetAll()
	if err != nil {
		log.Printf("failed to get battery info from system: %v", err)
		return
	}
	for i, b := range batts {
		n := mkId(i)
		pc := math.Abs(b.Current/b.Full) * 100.0
		self.Data[n] = append(self.Data[n], pc)
		self.Labels[n] = fmt.Sprintf("%3.0f%% %.0f/%.0f", pc, math.Abs(b.Current), math.Abs(b.Full))
	}
}
