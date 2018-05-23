// +build freebsd linux darwin

package widgets

import (
	"os/exec"
	"strconv"
	"strings"
)

func (self *Proc) update() {
	processes := Processes()
	// have to iterate like this in order to actually change the value
	for i, _ := range processes {
		processes[i].CPU /= self.cpuCount
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}

func Processes() []Process {
	output, _ := exec.Command("ps", "-axo", "pid,comm,pcpu,pmem,args").Output()
	// converts to []string and removes the header
	strOutput := strings.Split(strings.TrimSpace(string(output)), "\n")[1:]
	processes := []Process{}
	for _, line := range strOutput {
		split := strings.Fields(line)
		pid, _ := strconv.Atoi(split[0])
		cpu, _ := strconv.ParseFloat(split[2], 64)
		mem, _ := strconv.ParseFloat(split[3], 64)
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
