package widgets

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"

	"github.com/xxxserxxx/gotop/v4/devices"
	ui "github.com/xxxserxxx/gotop/v4/termui"
	"github.com/xxxserxxx/gotop/v4/utils"
)

type MemWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
}

func NewMemWidget(updateInterval time.Duration, horizontalScale int) *MemWidget {
	widg := &MemWidget{
		LineGraph:      ui.NewLineGraph(),
		updateInterval: updateInterval,
	}
	widg.Title = tr.Value("widget.label.mem")
	widg.HorizontalScale = horizontalScale
	mems := make(map[string]devices.MemoryInfo)
	devices.UpdateMem(mems)
	for name, mem := range mems {
		widg.Data[name] = []float64{0}
		widg.renderMemInfo(name, mem)
	}

	go func() {
		for range time.NewTicker(widg.updateInterval).C {
			widg.Lock()
			devices.UpdateMem(mems)
			for label, mi := range mems {
				widg.renderMemInfo(label, mi)
			}
			widg.Unlock()
		}
	}()

	return widg
}

func (mem *MemWidget) EnableMetric() {
	mems := make(map[string]devices.MemoryInfo)
	devices.UpdateMem(mems)
	for l, _ := range mems {
		lc := l
		metrics.NewGauge(makeName("memory", l), func() float64 {
			if ds, ok := mem.Data[lc]; ok {
				return ds[len(ds)-1]
			}
			return 0.0
		})
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
