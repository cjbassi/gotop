package devices

import (
	"log"
	"time"
)

var cpuFuncs []func(map[string]int, bool) map[string]error

// RegisterCPU adds a new CPU device to the CPU widget. labels returns the
// names of the devices; they should be as short as possible, and the indexes
// of the returned slice should align with the values returned by the percents
// function.  The percents function should return the percent CPU usage of the
// device(s), sliced over the time duration supplied.  If the bool argument to
// percents is true, it is expected that the return slice
//
// labels may be called once and the value cached.  This means the number of
// cores should not change dynamically.
func RegisterCPU(f func(map[string]int, bool) map[string]error) {
	cpuFuncs = append(cpuFuncs, f)
}

// CPUPercent calculates the percentage of cpu used either per CPU or combined.
// Returns one value per cpu, or a single value if percpu is set to false.
func UpdateCPU(cpus map[string]int, interval time.Duration, logical bool) {
	for _, f := range cpuFuncs {
		errs := f(cpus, logical)
		if errs != nil {
			for k, e := range errs {
				log.Printf("%s: %s", k, e)
			}
		}
	}
}
