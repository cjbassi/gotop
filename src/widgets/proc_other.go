// +build freebsd darwin

package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// Define column widths for ps output used in Processes()
	five  = "12345"
	ten   = five + five
	fifty = ten + ten + ten + ten + ten
)

func (self *Proc) update() {
	processes, err := Processes()
	if err != nil {
		log.Printf("failed to retrieve processes: %v", err)
		return
	}

	// have to iterate like this in order to actually change the value
	for i := range processes {
		processes[i].CPU /= self.cpuCount
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}

func Processes() ([]Process, error) {
	keywords := fmt.Sprintf("pid=%s,comm=%s,pcpu=%s,pmem=%s,args", ten, fifty, five, five)
	output, err := exec.Command("ps", "-wwcaxo", keywords).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ps' command: %v", err)
	}

	// converts to []string and removes the header
	strOutput := strings.Split(strings.TrimSpace(string(output)), "\n")[1:]

	processes := []Process{}
	for _, line := range strOutput {
		pid, err := strconv.Atoi(strings.TrimSpace(line[0:10]))
		if err != nil {
			log.Printf("failed to convert first field to int: %v. split: %v", err, line)
		}
		cpu, err := strconv.ParseFloat(strings.TrimSpace(line[63:68]), 64)
		if err != nil {
			log.Printf("failed to convert third field to float: %v. split: %v", err, line)
		}
		mem, err := strconv.ParseFloat(strings.TrimSpace(line[69:74]), 64)
		if err != nil {
			log.Printf("failed to convert fourth field to float: %v. split: %v", err, line)
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
	return processes, nil
}
