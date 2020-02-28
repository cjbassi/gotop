// +build !freebsd

package devices

import (
	psMem "github.com/shirou/gopsutil/mem"
)

func init() {
	mf := func(mems map[string]MemoryInfo) map[string]error {
		memory, err := psMem.SwapMemory()
		if err != nil {
			return map[string]error{"Swap": err}
		}
		mems["Swap"] = MemoryInfo{
			Total:       memory.Total,
			Used:        memory.Used,
			UsedPercent: memory.UsedPercent,
		}
		return nil
	}
	RegisterMem(mf)
}
