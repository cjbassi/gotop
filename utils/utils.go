package utils

import (
	"math"
)

func BytesToKB(b int) int {
	return int((float64(b) / math.Pow10(3)))
}

func BytesToMB(b int) int {
	return int((float64(b) / math.Pow10(6)))
}

func BytesToGB(b int) int {
	return int((float64(b) / math.Pow10(9)))
}
