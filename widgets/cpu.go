package widgets

import (
	"fmt"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/xxxserxxx/gotop/v4/devices"

	ui "github.com/xxxserxxx/gotop/v4/termui"
)

type CPUWidget struct {
	*ui.LineGraph
	CPUCount        int
	ShowAverageLoad bool
	ShowPerCPULoad  bool
	updateInterval  time.Duration
	updateLock      sync.Mutex
	cpuLoads        map[string]float64
}

var cpuLabels []string

func NewCPUWidget(updateInterval time.Duration, horizontalScale int, showAverageLoad bool, showPerCPULoad bool) *CPUWidget {
	self := &CPUWidget{
		LineGraph:       ui.NewLineGraph(),
		CPUCount:        len(cpuLabels),
		updateInterval:  updateInterval,
		ShowAverageLoad: showAverageLoad,
		ShowPerCPULoad:  showPerCPULoad,
		cpuLoads:        make(map[string]float64),
	}
	self.Title = " CPU Usage "
	self.HorizontalScale = horizontalScale

	if !(self.ShowAverageLoad || self.ShowPerCPULoad) {
		if self.CPUCount <= 8 {
			self.ShowPerCPULoad = true
		} else {
			self.ShowAverageLoad = true
		}
	}

	if self.ShowAverageLoad {
		self.Data["AVRG"] = []float64{0}
	}

	if self.ShowPerCPULoad {
		cpus := make(map[string]int)
		devices.UpdateCPU(cpus, self.updateInterval, self.ShowPerCPULoad)
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

const AVRG = "AVRG"

func (cpu *CPUWidget) EnableMetric() {
	if cpu.ShowAverageLoad {
		metrics.NewGauge(makeName("cpu", " avg"), func() float64 {
			return cpu.cpuLoads[AVRG]
		})
	} else {
		cpus := make(map[string]int)
		devices.UpdateCPU(cpus, cpu.updateInterval, cpu.ShowPerCPULoad)
		for key, perc := range cpus {
			kc := key
			cpu.cpuLoads[key] = float64(perc)
			metrics.NewGauge(makeName("cpu", key), func() float64 {
				return cpu.cpuLoads[kc]
			})
		}
	}
}

func (cpu *CPUWidget) Scale(i int) {
	cpu.LineGraph.HorizontalScale = i
}

func (cpu *CPUWidget) update() {
	if cpu.ShowAverageLoad {
		go func() {
			cpus := make(map[string]int)
			devices.UpdateCPU(cpus, cpu.updateInterval, false)
			cpu.Lock()
			defer cpu.Unlock()
			cpu.updateLock.Lock()
			defer cpu.updateLock.Unlock()
			var val float64
			for _, v := range cpus {
				val = float64(v)
				break
			}
			cpu.Data[AVRG] = append(cpu.Data[AVRG], val)
			cpu.Labels[AVRG] = fmt.Sprintf("%3.0f%%", val)
			cpu.cpuLoads[AVRG] = val
		}()
	}

	if cpu.ShowPerCPULoad {
		go func() {
			cpus := make(map[string]int)
			devices.UpdateCPU(cpus, cpu.updateInterval, true)
			cpu.Lock()
			defer cpu.Unlock()
			cpu.updateLock.Lock()
			defer cpu.updateLock.Unlock()
			for key, percent := range cpus {
				cpu.Data[key] = append(cpu.Data[key], float64(percent))
				cpu.Labels[key] = fmt.Sprintf("%d%%", percent)
				cpu.cpuLoads[key] = float64(percent)
			}
		}()
	}
}
