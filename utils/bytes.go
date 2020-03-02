package utils

import (
	"math"
)

var (
	KB = uint64(math.Pow(2, 10))
	MB = uint64(math.Pow(2, 20))
	GB = uint64(math.Pow(2, 30))
	TB = uint64(math.Pow(2, 40))
)

func CelsiusToFahrenheit(c int) int {
	return c*9/5 + 32
}

func BytesToKB(b uint64) float64 {
	return float64(b) / float64(KB)
}

func BytesToMB(b uint64) float64 {
	return float64(b) / float64(MB)
}

func BytesToGB(b uint64) float64 {
	return float64(b) / float64(GB)
}

func BytesToTB(b uint64) float64 {
	return float64(b) / float64(TB)
}

func ConvertBytes(b uint64) (float64, string) {
	switch {
	case b < KB:
		return float64(b), "B"
	case b < MB:
		return BytesToKB(b), "KB"
	case b < GB:
		return BytesToMB(b), "MB"
	case b < TB:
		return BytesToGB(b), "GB"
	default:
		return BytesToTB(b), "TB"
	}
}
