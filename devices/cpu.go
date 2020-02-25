package devices

import (
	"time"
)

var deviceCounts []func(bool) (int, error)
var devicePercents []func(time.Duration, bool) ([]float64, error)
var numDevices int

// Counts returns the number of CPUs registered.
//
// logical tells Counts to count the logical cores; this may be ignored for
// some devices.
func Counts(logical bool) (int, error) {
	var rv int
	var re error
	for _, d := range deviceCounts {
		r, err := d(logical)
		if err != nil {
			return rv, re
		}
		rv += r
	}
	return rv, re
}

// Percent calculates the percentage of cpu used either per CPU or combined.
// Returns one value per cpu, or a single value if percpu is set to false.
func Percent(interval time.Duration, combined bool) ([]float64, error) {
	var rvs []float64
	rvs = make([]float64, 0, numDevices)
	for _, f := range devicePercents {
		vs, err := f(interval, combined)
		if err != nil {
			return rvs, err
		}
		for _, v := range vs {
			rvs = append(rvs, v)
		}
	}
	numDevices = len(rvs)
	return rvs, nil
}
