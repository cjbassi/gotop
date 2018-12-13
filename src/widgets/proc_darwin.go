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
	output, err := exec.Command("ps", "-axo", "pid,comm,pcpu,pmem,args").Output()
	if err != nil {
		log.Printf("failed to execute 'ps' command: %v", err)
	}
	// converts to []string and removes the header
	strOutput := strings.Split(strings.TrimSpace(string(output)), "\n")[1:]
	processes := []Process{}
	for _, line := range strOutput {
		split := strings.Fields(line)
		pid, err := strconv.Atoi(split[0])
		if err != nil {
			log.Printf("failed to convert first field to int: %v. split: %v", err, split)
		}
		cpu, err := strconv.ParseFloat(split[2], 64)
		if err != nil {
			log.Printf("failed to convert third field to float: %v. split: %v", err, split)
		}
		mem, err := strconv.ParseFloat(split[3], 64)
		if err != nil {
			log.Printf("failed to convert fourth field to float: %v. split: %v", err, split)
		}
		process := Process{
			PID:     pid,
			Command: split[1],
			CPU:     cpu,
			Mem:     mem,
			Args:    strings.Join(split[4:], " "),
		}
		processes = append(processes, process)
	}
	return processes
}
