// +build !freebsd

package widgets

import (
	"log"

	psMem "github.com/shirou/gopsutil/mem"
)

func (self *MemWidget) updateSwapMemory() {
	swapMemory, err := psMem.SwapMemory()
	if err != nil {
		log.Printf("failed to get swap memory info from gopsutil: %v", err)
	} else {
		self.renderMemInfo("Swap", MemoryInfo{
			Total:       swapMemory.Total,
			Used:        swapMemory.Used,
			UsedPercent: swapMemory.UsedPercent,
		})
	}
}
