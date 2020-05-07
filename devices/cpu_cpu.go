package devices

import (
	"fmt"
	"time"

	psCpu "github.com/shirou/gopsutil/cpu"
)

// FIXME: broken % under Linux.  Doesn't reflect reality *at all*.
// FIXME: gotop CPU use high -- gopsutils again? Try rolling back.
func init() {
	f := func(cpus map[string]int, iv time.Duration, l bool) map[string]error {
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
