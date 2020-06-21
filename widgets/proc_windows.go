package widgets

import (
	"fmt"
	"log"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func getProcs() ([]Proc, error) {
	psProcs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf(tr.Value("widget.proc.err.gopsutil", err.Error()))
	}

	procs := make([]Proc, len(psProcs))
	for i, psProc := range psProcs {
		pid := psProc.Pid
		command, err := psProc.Name()
		if err != nil {
			sps := fmt.Sprintf("%v", psProc)
			si := strconv.Itoa(i)
			spid := fmt.Sprintf("%d", pid)
			log.Println(tr.Value("widget.proc.err.getcmd", err.Error(), sps, si, spid))
		}
		cpu, err := psProc.CPUPercent()
		if err != nil {
			sps := fmt.Sprintf("%v", psProc)
			si := strconv.Itoa(i)
			spid := fmt.Sprintf("%d", pid)
			log.Println(tr.Value("widget.proc.err.cpupercent", err.Error(), sps, si, spid))
		}
		mem, err := psProc.MemoryPercent()
		if err != nil {
			sps := fmt.Sprintf("%v", psProc)
			si := strconv.Itoa(i)
			spid := fmt.Sprintf("%d", pid)
			log.Println(tr.Value("widget.proc.err.mempercent", err.Error(), sps, si, spid))
		}

		procs[i] = Proc{
			Pid:         int(pid),
			CommandName: command,
			CPU:         cpu,
			Mem:         float64(mem),
			// getting command args using gopsutil's Cmdline and CmdlineSlice wasn't
			// working the last time I tried it, so we're just reusing 'command'
			FullCommand: command,
		}
	}

	return procs, nil
}
