package widgets

import (
	"fmt"
	"log"

	psProc "github.com/shirou/gopsutil/process"
)

func getProcs() ([]Proc, error) {
	psProcs, err := psProc.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes from gopsutil: %v", err)
	}

	procs := make([]Proc, len(psProcs))
	for i, psProc := range psProcs {
		pid := psProc.Pid
		command, err := psProc.Name()
		if err != nil {
			log.Printf("failed to get process command from gopsutil: %v. psProc: %v. i: %v. pid: %v", err, psProc, i, pid)
		}
		cpu, err := psProc.CPUPercent()
		if err != nil {
			log.Printf("failed to get process cpu usage from gopsutil: %v. psProc: %v. i: %v. pid: %v", err, psProc, i, pid)
		}
		mem, err := psProc.MemoryPercent()
		if err != nil {
			log.Printf("failed to get process memeory usage from gopsutil: %v. psProc: %v. i: %v. pid: %v", err, psProc, i, pid)
		}

		procs[i] = Proc{
			Pid:         int(pid),
			CommandName: command,
			Cpu:         cpu,
			Mem:         float64(mem),
			// getting command args using gopsutil's Cmdline and CmdlineSlice wasn't
			// working the last time I tried it, so we're just reusing 'command'
			FullCommand: command,
		}
	}

	return procs, nil
}
