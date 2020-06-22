package widgets

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/xxxserxxx/gotop/v4/utils"
)

type processList struct {
	ProcessInformation struct {
		Process []struct {
			Pid  string `json:"pid"`
			Comm string `json:"command"`
			CPU  string `json:"percent-cpu" `
			Mem  string `json:"percent-memory" `
			Args string `json:"arguments" `
		} `json:"process"`
	} `json:"process-information"`
}

func getProcs() ([]Proc, error) {
	output, err := exec.Command("ps", "-axo pid,comm,%cpu,%mem,args", "--libxo", "json").Output()
	if err != nil {
		return nil, fmt.Errorf(tr.Value("widget.proc.err.ps", err.Error()))
	}

	list := processList{}
	err = json.Unmarshal(output, &list)
	if err != nil {
		return nil, fmt.Errorf(tr.Value("widget.proc.err.parse", err.Error()))
	}
	procs := []Proc{}

	for _, process := range list.ProcessInformation.Process {
		if process.Comm == "idle" {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(process.Pid))
		if err != nil {
			sp := fmt.Sprintf("%v", process)
			log.Printf(tr.Value("widget.proc.err.pidconv", err.Error(), sp))
		}
		cpu, err := strconv.ParseFloat(utils.ConvertLocalizedString(process.CPU), 32)
		if err != nil {
			sp := fmt.Sprintf("%v", process)
			log.Printf(tr.Value("widget.proc.err.cpuconv", err.Error(), sp))
		}
		mem, err := strconv.ParseFloat(utils.ConvertLocalizedString(process.Mem), 32)
		if err != nil {
			sp := fmt.Sprintf("%v", process)
			log.Printf(tr.Value("widget.proc.err.memconv", err.Error(), sp))
		}
		proc := Proc{
			Pid:         pid,
			CommandName: process.Comm,
			CPU:         cpu,
			Mem:         mem,
			FullCommand: process.Args,
		}
		procs = append(procs, proc)
	}

	return procs, nil
}
