package utils

import (
	"math"
)

func BytesToKB(b uint64) float64 {
	return float64(b) / math.Pow10(3)
}

func BytesToMB(b uint64) float64 {
	return float64(b) / math.Pow10(6)
}

func BytesToGB(b uint64) float64 {
	return float64(b) / math.Pow10(9)
}
