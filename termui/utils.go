package termui

import (
	"math"
)

const DOTS = '…'

// MaxString trims a string and adds dots if its length is greater than l
func MaxString(s string, l int) string {
	if l <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) > l {
		r = r[:l]
		r[l-1] = DOTS
	}
	return string(r)
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
