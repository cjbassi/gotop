package widgets

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/xxxserxxx/gotop/v4/devices"
	ui "github.com/xxxserxxx/gotop/v4/termui"
	"github.com/xxxserxxx/gotop/v4/utils"
)

type MemWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
	metrics        map[string]prometheus.Gauge
}

func NewMemWidget(updateInterval time.Duration, horizontalScale int) *MemWidget {
	self := &MemWidget{
		LineGraph:      ui.NewLineGraph(),
		updateInterval: updateInterval,
	}
	self.Title = " Memory Usage "
	self.HorizontalScale = horizontalScale
	mems := make(map[string]devices.MemoryInfo)
	devices.UpdateMem(mems)
	for name, mem := range mems {
		self.Data[name] = []float64{0}
		self.renderMemInfo(name, mem)
	}

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			devices.UpdateMem(mems)
			for label, mi := range mems {
				self.renderMemInfo(label, mi)
				if self.metrics != nil && self.metrics[label] != nil {
					self.metrics[label].Set(mi.UsedPercent)
				}
			}
			self.Unlock()
		}
	}()

	return self
}

func (mem *MemWidget) EnableMetric() {
	mem.metrics = make(map[string]prometheus.Gauge)
	mems := make(map[string]devices.MemoryInfo)
	devices.UpdateMem(mems)
	for l, m := range mems {
		mem.metrics[l] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gotop",
			Subsystem: "memory",
			Name:      l,
		})
		mem.metrics[l].Set(m.UsedPercent)
		prometheus.MustRegister(mem.metrics[l])
	}
}

func (mem *MemWidget) Scale(i int) {
	mem.LineGraph.HorizontalScale = i
}

func (mem *MemWidget) renderMemInfo(line string, memoryInfo devices.MemoryInfo) {
	mem.Data[line] = append(mem.Data[line], memoryInfo.UsedPercent)
	memoryTotalBytes, memoryTotalMagnitude := utils.ConvertBytes(memoryInfo.Total)
	memoryUsedBytes, memoryUsedMagnitude := utils.ConvertBytes(memoryInfo.Used)
	mem.Labels[line] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
		memoryInfo.UsedPercent,
		memoryUsedBytes,
		memoryUsedMagnitude,
		memoryTotalBytes,
		memoryTotalMagnitude,
	)
}
