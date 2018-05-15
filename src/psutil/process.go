package psutil

import (
	// "fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Process represents each process.
type Process struct {
	PID     int
	Command string
	CPU     float64
	Mem     float64
}

func Processes() []Process {
	output, _ := exec.Command("ps", "-acxo", "pid,comm,pcpu,pmem").Output()
	strOutput := strings.TrimSpace(string(output))
	processes := []Process{}
	for _, line := range strings.Split(strOutput, "\n")[1:] {
		split := strings.Fields(line)
		// fmt.Println(split)
		pid, _ := strconv.Atoi(split[0])
		cpu, _ := strconv.ParseFloat(split[2], 64)
		mem, _ := strconv.ParseFloat(split[3], 64)
		process := Process{
			PID:     pid,
			Command: split[1],
			CPU:     cpu,
			Mem:     mem,
		}
		processes = append(processes, process)
	}
	return processes
}
