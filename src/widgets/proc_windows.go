package widgets

import (
	"log"

	psProc "github.com/shirou/gopsutil/process"
)

func (self *Proc) update() {
	psProcesses, err := psProc.Processes()
	if err != nil {
		log.Printf("failed to get processes from gopsutil: %v", err)
	}
	processes := make([]Process, len(psProcesses))
	for i, psProcess := range psProcesses {
		pid := psProcess.Pid
		command, err := psProcess.Name()
		if err != nil {
			log.Printf("failed to get process command from gopsutil: %v. psProcess: %v. i: %v. pid: %v", err, psProcess, i, pid)
		}
		cpu, err := psProcess.CPUPercent()
		if err != nil {
			log.Printf("failed to get process cpu usage from gopsutil: %v. psProcess: %v. i: %v. pid: %v", err, psProcess, i, pid)
		}
		mem, err := psProcess.MemoryPercent()
		if err != nil {
			log.Printf("failed to get process memeory usage from gopsutil: %v. psProcess: %v. i: %v. pid: %v", err, psProcess, i, pid)
		}

		processes[i] = Process{
			int(pid),
			command,
			cpu / self.cpuCount,
			float64(mem),
			// getting command args using gopsutil's Cmdline and CmdlineSlice wasn't
			// working the last time I tried it, so we're just reusing 'command'
			command,
		}
	}

	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
}
