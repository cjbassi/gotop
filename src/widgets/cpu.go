package widgets

import (
	"fmt"
	"log"
	"sync"
	"time"

	ui "github.com/cjbassi/gotop/src/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	Count        int  // number of cores
	Average      bool // show average load
	PerCPU       bool // show per-core load
	interval     time.Duration
	formatString string
	renderLock   *sync.RWMutex
}

func NewCPU(renderLock *sync.RWMutex, interval time.Duration, horizontalScale int, average bool, percpu bool) *CPU {
	count, err := psCPU.Counts(false)
	if err != nil {
		log.Printf("failed to get CPU count from gopsutil: %v", err)
	}
	formatString := "CPU%1d"
	if count > 10 {
		formatString = "CPU%02d"
	}
	self := &CPU{
		LineGraph:    ui.NewLineGraph(),
		Count:        count,
		interval:     interval,
		Average:      average,
		PerCPU:       percpu,
		formatString: formatString,
		renderLock:   renderLock,
	}
	self.Title = " CPU Usage "
	self.HorizontalScale = horizontalScale

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
			k := fmt.Sprintf(formatString, i)
			self.Data[k] = []float64{0}
		}
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

func (self *CPU) update() {
	if self.Average {
		go func() {
			percent, err := psCPU.Percent(self.interval, false)
			self.renderLock.RLock()
			defer self.renderLock.RUnlock()
			if err != nil {
				log.Printf("failed to get average CPU usage percent from gopsutil: %v. self.interval: %v. percpu: %v", err, self.interval, false)
			} else {
				self.Data["AVRG"] = append(self.Data["AVRG"], percent[0])
				self.Labels["AVRG"] = fmt.Sprintf("%3.0f%%", percent[0])
			}
		}()
	}

	if self.PerCPU {
		go func() {
			percents, err := psCPU.Percent(self.interval, true)
			self.renderLock.RLock()
			defer self.renderLock.RUnlock()
			if err != nil {
				log.Printf("failed to get CPU usage percents from gopsutil: %v. self.interval: %v. percpu: %v", err, self.interval, true)
			} else {
				if len(percents) != self.Count {
					log.Printf("error: number of CPU usage percents from gopsutil doesn't match CPU count. percents: %v. self.Count: %v", percents, self.Count)
				} else {
					for i, percent := range percents {
						k := fmt.Sprintf(self.formatString, i)
						self.Data[k] = append(self.Data[k], percent)
						self.Labels[k] = fmt.Sprintf("%3.0f%%", percent)
					}
				}
			}
		}()
	}
}
