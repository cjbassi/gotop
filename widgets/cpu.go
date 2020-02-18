package widgets

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	psCpu "github.com/shirou/gopsutil/cpu"

	ui "github.com/xxxserxxx/gotop/termui"
)

type CpuWidget struct {
	*ui.LineGraph
	CpuCount        int
	ShowAverageLoad bool
	ShowPerCpuLoad  bool
	updateInterval  time.Duration
	formatString    string
	updateLock      sync.Mutex
	metric          []prometheus.Gauge
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

func (self *CpuWidget) EnableMetric() {
	if self.ShowAverageLoad {
		self.metric = make([]prometheus.Gauge, 1)
		self.metric[0] = prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "cpu",
			Name:      "avg",
		})
	} else {
		ctx, ccl := context.WithTimeout(context.Background(), time.Second*5)
		defer ccl()
		percents, err := psCpu.PercentWithContext(ctx, self.updateInterval, true)
		if err != nil {
			log.Printf("error setting up metrics: %v", err)
			return
		}
		self.metric = make([]prometheus.Gauge, self.CpuCount)
		for i, perc := range percents {
			gauge := prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gotop",
				Subsystem: "cpu",
				Name:      fmt.Sprintf("%d", i),
			})
			gauge.Set(perc)
			prometheus.MustRegister(gauge)
			self.metric[i] = gauge
		}
	}
}

func (b *CpuWidget) Scale(i int) {
	b.LineGraph.HorizontalScale = i
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
				if self.metric != nil {
					self.metric[0].Set(percent[0])
				}
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
						if self.metric != nil {
							if self.metric[i] == nil {
								log.Printf("ERROR: not enough metrics %d", i)
							} else {
								self.metric[i].Set(percent)
							}
						}
					}
				}
			}
		}()
	}
}
