package widgets

import (
	"fmt"
	"log"
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

	go func() {
		ticker := time.NewTicker(self.interval)
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Mem) update() {
	main, err := psMem.VirtualMemory()
	if err != nil {
		log.Printf("failed to get main memory info from gopsutil: %v", err)
	} else {
		self.Data["Main"] = append(self.Data["Main"], main.UsedPercent)
		mainTotalBytes, mainTotalMagnitude := utils.ConvertBytes(main.Total)
		mainUsedBytes, mainUsedMagnitude := utils.ConvertBytes(main.Used)
		self.Labels["Main"] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s", main.UsedPercent, mainUsedBytes, mainUsedMagnitude, mainTotalBytes, mainTotalMagnitude)
	}

	swap, err := psMem.SwapMemory()
	if err != nil {
		log.Printf("failed to get swap memory info from gopsutil: %v", err)
	} else {
		self.Data["Swap"] = append(self.Data["Swap"], swap.UsedPercent)
		swapTotalBytes, swapTotalMagnitude := utils.ConvertBytes(swap.Total)
		swapUsedBytes, swapUsedMagnitude := utils.ConvertBytes(swap.Used)
		self.Labels["Swap"] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s", swap.UsedPercent, swapUsedBytes, swapUsedMagnitude, swapTotalBytes, swapTotalMagnitude)
	}
}
