package widgets

import (
	"strconv"
	"time"

	ui "github.com/cjbassi/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	Count    int // number of CPUs
	interval time.Duration
}

func NewCPU(interval time.Duration, zoom int) *CPU {
	count, _ := psCPU.Counts(false)
	self := &CPU{
		LineGraph: ui.NewLineGraph(),
		Count:     count,
		interval:  interval,
	}
	self.Label = "CPU Usage"
	self.Zoom = zoom
	if self.Count <= 8 {
		for i := 0; i < self.Count; i++ {
			key := "CPU" + strconv.Itoa(i+1)
			self.Data[key] = []float64{0}
		}
	} else {
		self.Data["Average"] = []float64{0}
	}

	go self.update()
	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *CPU) update() {
	// psutil calculates the CPU usage over a 1 second interval, therefore it blocks for 1 second
	if self.Count <= 8 {
		percent, _ := psCPU.Percent(self.interval, true)
		for i := 0; i < self.Count; i++ {
			key := "CPU" + strconv.Itoa(i+1)
			self.Data[key] = append(self.Data[key], percent[i])
		}
	} else {
		percent, _ := psCPU.Percent(self.interval, false)
		self.Data["Average"] = append(self.Data["Average"], percent[0])
	}
}
