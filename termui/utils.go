package termui

import (
	"math"
)

const DOTS = '…'

// MaxString trims a string and adds dots if the string is longer than a give length.
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
