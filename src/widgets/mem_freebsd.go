package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

type SwapInfo struct {
	Total      uint64
	Used       uint64
	Percentage float64
}

func convert(s []string) (SwapInfo, error) {
	total, err := strconv.ParseUint(s[0], 10, 64)
	if err != nil {
		return SwapInfo{}, fmt.Errorf("int converion failed %v", err)
	}

	used, err := strconv.ParseUint(s[1], 10, 64)
	if err != nil {
		return SwapInfo{}, fmt.Errorf("int converion failed %v", err)
	}

	percentage, err := strconv.ParseFloat(s[2], 64)
	if err != nil {
		return SwapInfo{}, fmt.Errorf("float converion failed %v", err)
	}

	return SwapInfo{
		Total:      total * utils.KB,
		Used:       used * utils.KB,
		Percentage: percentage,
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

//TODO This can be DRYed more
func (self *MemWidget) updateSwapMemory() {
	swapMemory, err := GatherSwapInfo()
	if err != nil {
		log.Printf("failed to get swap memory info from gopsutil: %v", err)
	} else {
		self.Data["Swap"] = append(self.Data["Swap"], swapMemory.Percentage)
		swapMemoryTotalBytes, swapMemoryTotalMagnitude := utils.ConvertBytes(swapMemory.Total)
		swapMemoryUsedBytes, swapMemoryUsedMagnitude := utils.ConvertBytes(swapMemory.Used)
		self.Labels["Swap"] = fmt.Sprintf("%3.0f%% %5.1f%s/%.0f%s",
			swapMemory.Percentage,
			swapMemoryUsedBytes,
			swapMemoryUsedMagnitude,
			swapMemoryTotalBytes,
			swapMemoryTotalMagnitude,
		)
	}
}
