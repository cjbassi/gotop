package widgets

import (
	"fmt"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	Count    int  // number of cores
	Average  bool // show average load
	PerCPU   bool // show per-core load
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

	if !(self.Average || self.PerCPU) {
		if self.Count <= 8 {
			self.PerCPU = true
		} else {
			self.Average = true
		}
	}

	if self.Average {
		self.Data["AVRG"] = []float64{0}
	}

	if self.PerCPU {
		for i := 0; i < self.Count; i++ {
			k := fmt.Sprintf("CPU%d", i)
			self.Data[k] = []float64{0}
		}
	}

	go self.update() // update asynchronously because of 1 second blocking period

	go func() {
		ticker := time.NewTicker(self.interval)
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

// calculates the CPU usage over a 1 second interval and blocks for the duration
func (self *CPU) update() {
	if self.Average {
		go func() {
			percent, _ := psCPU.Percent(self.interval, false)
			self.Data["AVRG"] = append(self.Data["AVRG"], percent[0])
			self.Labels["AVRG"] = fmt.Sprintf("%3.0f%%", percent[0])
		}()
	}

	if self.PerCPU {
		go func() {
			percents, _ := psCPU.Percent(self.interval, true)
			if len(percents) != self.Count {
				count, _ := psCPU.Counts(false)
				utils.Error("CPU percentages",
					fmt.Sprint(
						"self.Count: ", self.Count, "\n",
						"gopsutil.Counts(): ", count, "\n",
						"len(percents): ", len(percents), "\n",
						"percents: ", percents, "\n",
						"self.interval: ", self.interval,
					))
			}
			for i := 0; i < self.Count; i++ {
				k := fmt.Sprintf("CPU%d", i)
				self.Data[k] = append(self.Data[k], percents[i])
				self.Labels[k] = fmt.Sprintf("%3.0f%%", percents[i])
			}
		}()
	}
}
