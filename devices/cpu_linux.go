// +build linux

package devices

import "github.com/shirou/gopsutil/cpu"

func CpuCount() (int, error) {
	cpuCount, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}
	if cpuCount == 0 {
		is, err := cpu.Info()
		if err != nil {
			return 0, err
		}
		if is[0].Cores > 0 {
			return len(is) / 2, nil
		}
		return len(is), nil
	}
	return cpuCount, nil
}
