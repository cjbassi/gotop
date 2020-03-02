package widgets

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xxxserxxx/gotop/v3/devices"

	ui "github.com/xxxserxxx/gotop/v3/termui"
)

type CpuWidget struct {
	*ui.LineGraph
	CpuCount        int
	ShowAverageLoad bool
	ShowPerCpuLoad  bool
	updateInterval  time.Duration
	updateLock      sync.Mutex
	metric          map[string]prometheus.Gauge
}

var cpuLabels []string

func NewCpuWidget(updateInterval time.Duration, horizontalScale int, showAverageLoad bool, showPerCpuLoad bool) *CpuWidget {
	self := &CpuWidget{
		LineGraph:       ui.NewLineGraph(),
		CpuCount:        len(cpuLabels),
		updateInterval:  updateInterval,
		ShowAverageLoad: showAverageLoad,
		ShowPerCpuLoad:  showPerCpuLoad,
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
		cpus := make(map[string]int)
		devices.UpdateCPU(cpus, self.updateInterval, self.ShowPerCpuLoad)
		for k, v := range cpus {
			self.Data[k] = []float64{float64(v)}
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
		self.metric = make(map[string]prometheus.Gauge)
		self.metric["AVRG"] = prometheus.NewGauge(prometheus.GaugeOpts{
			Subsystem: "cpu",
			Name:      "avg",
		})
	} else {
		cpus := make(map[string]int)
		devices.UpdateCPU(cpus, self.updateInterval, self.ShowPerCpuLoad)
		self.metric = make(map[string]prometheus.Gauge)
		for key, perc := range cpus {
			gauge := prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "gotop",
				Subsystem: "cpu",
				Name:      key,
			})
			gauge.Set(float64(perc))
			prometheus.MustRegister(gauge)
			self.metric[key] = gauge
		}
	}
}

func (b *CpuWidget) Scale(i int) {
	b.LineGraph.HorizontalScale = i
}

func (self *CpuWidget) update() {
	if self.ShowAverageLoad {
		go func() {
			cpus := make(map[string]int)
			devices.UpdateCPU(cpus, self.updateInterval, false)
			self.Lock()
			defer self.Unlock()
			self.updateLock.Lock()
			defer self.updateLock.Unlock()
			var val float64
			for _, v := range cpus {
				val = float64(v)
				break
			}
			self.Data["AVRG"] = append(self.Data["AVRG"], val)
			self.Labels["AVRG"] = fmt.Sprintf("%3.0f%%", val)
			if self.metric != nil {
				self.metric["AVRG"].Set(val)
			}
		}()
	}

	if self.ShowPerCpuLoad {
		go func() {
			cpus := make(map[string]int)
			devices.UpdateCPU(cpus, self.updateInterval, true)
			self.Lock()
			defer self.Unlock()
			self.updateLock.Lock()
			defer self.updateLock.Unlock()
			for key, percent := range cpus {
				self.Data[key] = append(self.Data[key], float64(percent))
				self.Labels[key] = fmt.Sprintf("%d%%", percent)
				if self.metric != nil {
					if self.metric[key] == nil {
						log.Printf("no metrics for %s", key)
					} else {
						self.metric[key].Set(float64(percent))
					}
				}
			}
		}()
	}
}
