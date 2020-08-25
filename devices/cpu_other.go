// +build !linux

package devices

import "github.com/shirou/gopsutil/cpu"

func CpuCount() (int, error) {
	return cpu.Counts(false)
}
