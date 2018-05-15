package widgets

import (
	psProc "github.com/shirou/gopsutil/process"
)

func (self *Proc) update() {
	psProcesses, _ := psProc.Processes()
	processes := make([]Process, len(psProcesses))
	for i, psProcess := range psProcesses {
		pid := psProcess.Pid
		command, _ := psProcess.Name()
		cpu, _ := psProcess.CPUPercent()
		mem, _ := psProcess.MemoryPercent()

		processes[i] = Process{
			int(pid),
			command,
			cpu / self.cpuCount,
			float64(mem),
		}
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}
