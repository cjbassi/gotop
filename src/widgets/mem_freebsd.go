package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

func convert(s []string) (MemoryInfo, error) {
	total, err := strconv.ParseUint(s[0], 10, 64)
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("int converion failed %v", err)
	}

	used, err := strconv.ParseUint(s[1], 10, 64)
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("int converion failed %v", err)
	}

	percentage, err := strconv.ParseFloat(strings.TrimSuffix(s[2], "%"), 64)
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("float converion failed %v", err)
	}

	return MemoryInfo{
		Total:       total * utils.KB,
		Used:        used * utils.KB,
		UsedPercent: percentage,
	}, nil
}

func gatherSwapInfo() (MemoryInfo, error) {
	cmd := "swapinfo -k|sed -n '1!p'|awk '{print $2,$3,$5}'"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		if err != nil {
			return MemoryInfo{}, fmt.Errorf("command failed %v", err)
		}
	}

	ss := strings.Split(strings.TrimSuffix(string(output), "\n"), " ")

	return convert(ss)
}

func (self *MemWidget) updateSwapMemory() {
	swapMemory, err := gatherSwapInfo()
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
