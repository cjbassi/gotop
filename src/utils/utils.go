package utils

import (
	"fmt"
	"math"

	ui "github.com/cjbassi/termui"
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

func ConvertBytes(b uint64) (float64, string) {
	if b >= 1000000000 {
		return BytesToGB(uint64(b)), "GB"
	} else if b >= 1000000 {
		return BytesToMB(uint64(b)), "MB"
	} else if b >= 1000 {
		return BytesToKB(uint64(b)), "KB"
	} else {
		return float64(b), "B"
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Error(issue, diagnostics string) {
	ui.Close()
	fmt.Println("Error caught. Exiting program.")
	fmt.Println()
	fmt.Println("Issue with " + issue + ".")
	fmt.Println()
	fmt.Println("Diagnostics:\n" + diagnostics)
	fmt.Println()
	panic(1)
}
