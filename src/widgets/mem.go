package widgets

import (
	"fmt"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psMem "github.com/shirou/gopsutil/mem"
)

type Mem struct {
	*ui.LineGraph
	interval time.Duration
}

func NewMem(interval time.Duration, zoom int) *Mem {
	self := &Mem{
		LineGraph: ui.NewLineGraph(),
		interval:  interval,
	}
	self.Label = "Memory Usage"
	self.Zoom = zoom
	self.Data["Main"] = []float64{0}
	self.Data["Swap"] = []float64{0}

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Mem) update() {
	main, _ := psMem.VirtualMemory()
	swap, _ := psMem.SwapMemory()
	self.Data["Main"] = append(self.Data["Main"], main.UsedPercent)
	self.Data["Swap"] = append(self.Data["Swap"], swap.UsedPercent)

	mainTotalBytes, mainTotalMagnitude := utils.ConvertBytes(main.Total)
	swapTotalBytes, swapTotalMagnitude := utils.ConvertBytes(swap.Total)
	mainUsedBytes, mainUsedMagnitude := utils.ConvertBytes(main.Used)
	swapUsedBytes, swapUsedMagnitude := utils.ConvertBytes(swap.Used)
	self.Labels["Main"] = fmt.Sprintf("%3.0f%% %.0f%s/%.0f%s", main.UsedPercent, mainUsedBytes, mainUsedMagnitude, mainTotalBytes, mainTotalMagnitude)
	self.Labels["Swap"] = fmt.Sprintf("%3.0f%% %.0f%s/%.0f%s", swap.UsedPercent, swapUsedBytes, swapUsedMagnitude, swapTotalBytes, swapTotalMagnitude)
}
