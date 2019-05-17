package widgets

import (
	"fmt"
	"log"

	psMem "github.com/shirou/gopsutil/mem"

	"github.com/cjbassi/gotop/src/utils"
)

func (self *MemWidget) update() {
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

	swapMemory, err := psMem.SwapMemory()
	if err != nil {
		log.Printf("failed to get swap memory info from gopsutil: %v", err)
	} else {
		self.Data["Swap"] = append(self.Data["Swap"], swapMemory.UsedPercent)
		swapMemoryTotalBytes, swapMemoryTotalMagnitude := utils.ConvertBytes(swapMemory.Total)
		swapMemoryUsedBytes, swapMemoryUsedMagnitude := utils.ConvertBytes(swapMemory.Used)
		self.Labels["Swap"] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
			swapMemory.UsedPercent,
			swapMemoryUsedBytes,
			swapMemoryUsedMagnitude,
			swapMemoryTotalBytes,
			swapMemoryTotalMagnitude,
		)
	}
}
