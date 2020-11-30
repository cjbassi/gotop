// +build linux

package devices

// TODO gopsutil is at v3, and we're using v2. See if v3 is released and upgrade if so.
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
