// +build linux freebsd

package widgets

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func (self *Proc) update() {
	processes := Processes()
	// have to iterate like this in order to actually change the value
	for i := range processes {
		processes[i].CPU /= self.cpuCount
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}

func Processes() []Process {
	output, err := exec.Command("ps", "-axo", "pid:10,comm:50,pcpu:5,pmem:5,args").Output()
	if err != nil {
		log.Printf("failed to execute 'ps' command: %v", err)
	}

	// converts to []string, removing trailing newline and header
	processStrArr := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")[1:]

	processes := []Process{}
	for _, line := range processStrArr {
		pid, err := strconv.Atoi(strings.TrimSpace(line[0:10]))
		if err != nil {
			log.Printf("failed to convert PID to int: %v. line: %v", err, line)
		}
		cpu, err := strconv.ParseFloat(strings.TrimSpace(line[63:68]), 64)
		if err != nil {
			log.Printf("failed to convert CPU usage to float: %v. line: %v", err, line)
		}
		mem, err := strconv.ParseFloat(strings.TrimSpace(line[69:74]), 64)
		if err != nil {
			log.Printf("failed to convert Mem usage to float: %v. line: %v", err, line)
		}
		process := Process{
			PID:     pid,
			Command: strings.TrimSpace(line[11:61]),
			CPU:     cpu,
			Mem:     mem,
			Args:    line[74:],
		}
		processes = append(processes, process)
	}
	return processes
}
