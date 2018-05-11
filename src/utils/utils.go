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
