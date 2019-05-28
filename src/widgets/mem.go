package widgets

import (
	"fmt"
	"log"
	"time"

	psMem "github.com/shirou/gopsutil/mem"

	ui "github.com/cjbassi/gotop/src/termui"
	"github.com/cjbassi/gotop/src/utils"
)

type MemWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
}

type MemoryInfo struct {
	Total       uint64
	Used        uint64
	UsedPercent float64
}

func (self *MemWidget) renderMemInfo(line string, memoryInfo MemoryInfo) {
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

func (self *MemWidget) updateMainMemory() {
	mainMemory, err := psMem.VirtualMemory()
	if err != nil {
		log.Printf("failed to get main memory info from gopsutil: %v", err)
	} else {
		self.renderMemInfo("Main", MemoryInfo{
			Total:       mainMemory.Total,
			Used:        mainMemory.Used,
			UsedPercent: mainMemory.UsedPercent,
		})
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
