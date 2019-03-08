package widgets

import (
	"fmt"
	"log"
	"sync"
	"time"

	psCpu "github.com/shirou/gopsutil/cpu"

	ui "github.com/cjbassi/gotop/src/termui"
)

type CpuWidget struct {
	*ui.LineGraph
	CpuCount        int
	ShowAverageLoad bool
	ShowPerCpuLoad  bool
	updateInterval  time.Duration
	formatString    string
	updateLock      sync.Mutex
}

func NewCpuWidget(updateInterval time.Duration, horizontalScale int, showAverageLoad bool, showPerCpuLoad bool) *CpuWidget {
	cpuCount, err := psCpu.Counts(false)
	if err != nil {
		log.Printf("failed to get CPU count from gopsutil: %v", err)
	}
	formatString := "CPU%1d"
	if cpuCount > 10 {
		formatString = "CPU%02d"
	}
	self := &CpuWidget{
		LineGraph:       ui.NewLineGraph(),
		CpuCount:        cpuCount,
		updateInterval:  updateInterval,
		ShowAverageLoad: showAverageLoad,
		ShowPerCpuLoad:  showPerCpuLoad,
		formatString:    formatString,
	}
	self.Title = " CPU Usage "
	self.HorizontalScale = horizontalScale

	if !(self.ShowAverageLoad || self.ShowPerCpuLoad) {
		if self.CpuCount <= 8 {
			self.ShowPerCpuLoad = true
		} else {
			self.ShowAverageLoad = true
		}
	}

	if self.ShowAverageLoad {
		self.Data["AVRG"] = []float64{0}
	}

	if self.ShowPerCpuLoad {
		for i := 0; i < int(self.CpuCount); i++ {
			key := fmt.Sprintf(formatString, i)
			self.Data[key] = []float64{0}
		}
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.update()
		}
	}()

	return self
}

func (self *CpuWidget) update() {
	if self.ShowAverageLoad {
		go func() {
			percent, err := psCpu.Percent(self.updateInterval, false)
			if err != nil {
				log.Printf("failed to get average CPU usage percent from gopsutil: %v. self.updateInterval: %v. percpu: %v", err, self.updateInterval, false)
			} else {
				self.Lock()
				defer self.Unlock()
				self.updateLock.Lock()
				defer self.updateLock.Unlock()
				self.Data["AVRG"] = append(self.Data["AVRG"], percent[0])
				self.Labels["AVRG"] = fmt.Sprintf("%3.0f%%", percent[0])
			}
		}()
	}

	if self.ShowPerCpuLoad {
		go func() {
			percents, err := psCpu.Percent(self.updateInterval, true)
			if err != nil {
				log.Printf("failed to get CPU usage percents from gopsutil: %v. self.updateInterval: %v. percpu: %v", err, self.updateInterval, true)
			} else {
				if len(percents) != int(self.CpuCount) {
					log.Printf("error: number of CPU usage percents from gopsutil doesn't match CPU count. percents: %v. self.Count: %v", percents, self.CpuCount)
				} else {
					self.Lock()
					defer self.Unlock()
					self.updateLock.Lock()
					defer self.updateLock.Unlock()
					for i, percent := range percents {
						key := fmt.Sprintf(self.formatString, i)
						self.Data[key] = append(self.Data[key], percent)
						self.Labels[key] = fmt.Sprintf("%3.0f%%", percent)
					}
				}
			}
		}()
	}
}
