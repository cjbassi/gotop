package widgets

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/src/termui"
	"github.com/distatus/battery"
)

type Batt struct {
	*ui.LineGraph
	Count    int // number of batteries
	interval time.Duration
}

func NewBatt(interval time.Duration, zoom int) *Batt {
	batts, err := battery.GetAll()
	self := &Batt{
		LineGraph: ui.NewLineGraph(),
		Count:     len(batts),
		interval:  interval,
	}
	self.Title = "Battery Status"
	self.Zoom = zoom
	if err != nil {
		log.Printf("failed to get battery info from system: %v", err)
	}
	for i, b := range batts {
		pc := math.Abs(b.Current/b.Full) * 100.0
		self.Data[mkId(i)] = []float64{pc}
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

func mkId(i int) string {
	return "Batt" + strconv.Itoa(i)
}

func (self *Batt) update() {
	batts, err := battery.GetAll()
	if err != nil {
		log.Printf("failed to get battery info from system: %v", err)
	}
	for i, b := range batts {
		n := mkId(i)
		pc := math.Abs(b.Current/b.Full) * 100.0
		self.Data[n] = append(self.Data[n], pc)
		self.Labels[n] = fmt.Sprintf("%3.0f%% %.0f/%.0f", pc, math.Abs(b.Current), math.Abs(b.Full))
	}
}
