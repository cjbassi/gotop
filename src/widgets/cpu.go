package widgets

import (
	"fmt"
	"time"

	ui "github.com/cjbassi/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	Count    int // number of cores
	Average bool // show average load
	PerCPU bool // show per-core load
	interval time.Duration
}

func NewCPU(interval time.Duration, zoom int, average bool, percpu bool) *CPU {
	count, _ := psCPU.Counts(false)
	self := &CPU{
		LineGraph: ui.NewLineGraph(),
		Count:     count,
		interval:  interval,
		Average:   average,
		PerCPU:    percpu,
	}
	self.Label = "CPU Usage"
	self.Zoom = zoom

	if self.Average {
		self.Data["Average"] = []float64{0}
	}

	if self.PerCPU {
		for i := 0; i < self.Count; i++ {
			k := fmt.Sprintf("CPU%d", i)
			self.Data[k] = []float64{0}
		}
	}

	ticker := time.NewTicker(self.interval)
	go func() {
		self.update()
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

// calculates the CPU usage over a 1 second interval and blocks for the duration
func (self *CPU) update() {
	if self.Average {
		percent, _ := psCPU.Percent(self.interval, false)
		self.Data["Average"] = append(self.Data["Average"], percent[0])
		self.Labels["Average"] = fmt.Sprintf("%3.0f%%", percent[0])
	}

	if self.PerCPU {
		percents, _ := psCPU.Percent(self.interval, true)
		for i := 0; i < self.Count; i++ {
			k := fmt.Sprintf("CPU%d", i)
			self.Data[k] = append(self.Data[k], percents[i])
			self.Labels[k] = fmt.Sprintf("%3.0f%%", percents[i])
		}
	}
}
