package devices

import (
	psCpu "github.com/shirou/gopsutil/cpu"
)

func init() {
	deviceCounts = append(deviceCounts, psCpu.Counts)
	devicePercents = append(devicePercents, psCpu.Percent)
}
