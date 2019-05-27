package widgets

import (
	"fmt"
	"github.com/cjbassi/gotop/src/utils"
	"log"
	"time"

	ui "github.com/cjbassi/gotop/src/termui"
	psMem "github.com/shirou/gopsutil/mem"
)

type MemWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
}

func (self *MemWidget) updateMainMemory() {
	mainMemory, err := psMem.VirtualMemory()
	if err != nil {
		log.Printf("failed to get main memory info from gopsutil: %v", err)
	} else {
		self.Data["Main"] = append(self.Data["Main"], mainMemory.UsedPercent)
		mainMemoryTotalBytes, mainMemoryTotalMagnitude := utils.ConvertBytes(mainMemory.Total)
		mainMemoryUsedBytes, mainMemoryUsedMagnitude := utils.ConvertBytes(mainMemory.Used)
		self.Labels["Main"] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
			mainMemory.UsedPercent,
			mainMemoryUsedBytes,
			mainMemoryUsedMagnitude,
			mainMemoryTotalBytes,
			mainMemoryTotalMagnitude,
		)
	}
}

func NewMemWidget(updateInterval time.Duration, horizontalScale int) *MemWidget {
	self := &MemWidget{
		LineGraph:      ui.NewLineGraph(),
		updateInterval: updateInterval,
	}
	self.Title = " Memory Usage "
	self.HorizontalScale = horizontalScale
	self.Data["Main"] = []float64{0}
	self.Data["Swap"] = []float64{0}

	self.updateMainMemory()
	self.updateSwapMemory()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.updateMainMemory()
			self.updateSwapMemory()
			self.Unlock()
		}
	}()

	return self
}
