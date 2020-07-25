package devices

import (
	"fmt"

	psCpu "github.com/shirou/gopsutil/cpu"
)

func init() {
	f := func(cpus map[string]int, l bool) map[string]error {
		cpuCount, err := psCpu.Counts(l)
		if err != nil {
			return nil
		}
		formatString := "CPU%1d"
		if cpuCount > 10 {
			formatString = "CPU%02d"
		}
		vals, err := psCpu.Percent(0, l)
		if err != nil {
			return map[string]error{"gopsutil": err}
		}
		for i := 0; i < len(vals); i++ {
			key := fmt.Sprintf(formatString, i)
			v := vals[i]
			if v > 100 {
				v = 100
			}
			cpus[key] = int(v)
		}
		return nil
	}
	RegisterCPU(f)
}
