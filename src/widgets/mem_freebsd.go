// +build freebsd
package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
	psMem "github.com/shirou/gopsutil/mem"
)

type SwapInfo struct {
	Total      uint64
	Used       uint64
	Percentage string
}

func convert(s []string) (SwapInfo, error) {
	total, err := strconv.ParseUint(s[0], 10, 64)
	if err != nil {
		return SwapInfo{}, fmt.Errorf("float converion failed %v", err)
	}

	used, err := strconv.ParseUint(s[1], 10, 64)
	if err != nil {
		return SwapInfo{}, fmt.Errorf("float converion failed %v", err)
	}
	return SwapInfo{
		Total:      total * utils.KB,
		Used:       used * utils.KB,
		Percentage: s[2],
	}, nil
}

func GatherSwapInfo() (SwapInfo, error) {
	cmd := "swapinfo -k|sed -n '1!p'|awk '{print $2,$3,$5}'"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		if err != nil {
			return SwapInfo{}, fmt.Errorf("command failed %v", err)
		}
	}

	ss := strings.Split(strings.TrimSuffix(string(output), "\n"), " ")

	return convert(ss)
}

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

	swapMemory, err := GatherSwapInfo()
	if err != nil {
		log.Printf("failed to get swap memory info from gopsutil: %v", err)
	} else {
		self.Data["Swap"] = append(self.Data["Swap"], 0)
		swapMemoryTotalBytes, swapMemoryTotalMagnitude := utils.ConvertBytes(swapMemory.Total)
		swapMemoryUsedBytes, swapMemoryUsedMagnitude := utils.ConvertBytes(swapMemory.Used)
		self.Labels["Swap"] = fmt.Sprintf("%3s   %5.1f%s/%.0f%s",
			swapMemory.Percentage,
			swapMemoryUsedBytes,
			swapMemoryUsedMagnitude,
			swapMemoryTotalBytes,
			swapMemoryTotalMagnitude,
		)
	}
}
