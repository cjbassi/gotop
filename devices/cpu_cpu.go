package devices

import (
	"fmt"
	"time"

	psCpu "github.com/shirou/gopsutil/cpu"
)

func init() {
	f := func(cpus map[string]float64, iv time.Duration, l bool) map[string]error {
		cpuCount, err := psCpu.Counts(l)
		if err != nil {
			return nil
		}
		formatString := "CPU%1d"
		if cpuCount > 10 {
			formatString = "CPU%02d"
		}
		vals, err := psCpu.Percent(iv, l)
		if err != nil {
			return map[string]error{"gopsutil": err}
		}
		for i := 0; i < len(vals); i++ {
			key := fmt.Sprintf(formatString, i)
			cpus[key] = vals[i]
		}
		return nil
	}
	RegisterCPU(f)
}
