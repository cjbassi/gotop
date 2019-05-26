package widgets

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

type processList struct {
	ProcessInformation struct {
		Process []struct {
			Pid  string `json:"pid"`
			Comm string `json:"command"`
			Cpu  string `json:"percent-cpu" `
			Mem  string `json:"percent-memory" `
			Args string `json:"arguments" `
		} `json:"process"`
	} `json:"process-information"`
}

func getProcs() ([]Proc, error) {
	output, err := exec.Command("ps", "-axo pid,comm,%cpu,%mem,args", "--libxo", "json").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ps' command: %v", err)
	}

	list := processList{}
	err = json.Unmarshal(output, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json. %s", err)
	}
	procs := []Proc{}

	for _, process := range list.ProcessInformation.Process {
		if process.Comm == "idle" {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(process.Pid))
		if err != nil {
			log.Printf("failed to convert first field to int: %v. split: %v", err, process)
		}
		cpu, err := strconv.ParseFloat(utils.ConvertLocalizedString(process.Cpu), 32)
		if err != nil {
			log.Printf("failed to convert third field to float: %v. split: %v", err, process)
		}
		mem, err := strconv.ParseFloat(utils.ConvertLocalizedString(process.Mem), 32)
		if err != nil {
			log.Printf("failed to convert fourth field to float: %v. split: %v", err, process)
		}
		proc := Proc{
			Pid:         pid,
			CommandName: process.Comm,
			Cpu:         cpu,
			Mem:         mem,
			FullCommand: process.Args,
		}
		procs = append(procs, proc)
	}

	return procs, nil
}
