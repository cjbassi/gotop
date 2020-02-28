package widgets

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/xxxserxxx/gotop/devices"
	ui "github.com/xxxserxxx/gotop/termui"
	"github.com/xxxserxxx/gotop/utils"
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
		self.Data[name] = []float64{mem.UsedPercent}
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

func (b *MemWidget) EnableMetric() {
	b.metrics = make(map[string]prometheus.Gauge)
	mems := make(map[string]devices.MemoryInfo)
	devices.UpdateMem(mems)
	for l, mem := range mems {
		b.metrics[l] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gotop",
			Subsystem: "memory",
			Name:      l,
		})
		b.metrics[l].Set(mem.UsedPercent)
		prometheus.MustRegister(b.metrics[l])
	}
}

func (b *MemWidget) Scale(i int) {
	b.LineGraph.HorizontalScale = i
}

func (self *MemWidget) renderMemInfo(line string, memoryInfo devices.MemoryInfo) {
	self.Data[line] = append(self.Data[line], memoryInfo.UsedPercent)
	memoryTotalBytes, memoryTotalMagnitude := utils.ConvertBytes(memoryInfo.Total)
	memoryUsedBytes, memoryUsedMagnitude := utils.ConvertBytes(memoryInfo.Used)
	self.Labels[line] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
		memoryInfo.UsedPercent,
		memoryUsedBytes,
		memoryUsedMagnitude,
		memoryTotalBytes,
		memoryTotalMagnitude,
	)
}
