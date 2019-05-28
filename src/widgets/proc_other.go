// +build darwin openbsd

package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

const (
	// Define column widths for ps output used in Procs()
	five  = "12345"
	ten   = five + five
	fifty = ten + ten + ten + ten + ten
)

func getProcs() ([]Proc, error) {
	keywords := fmt.Sprintf("pid=%s,comm=%s,pcpu=%s,pmem=%s,args", ten, fifty, five, five)
	output, err := exec.Command("ps", "-caxo", keywords).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ps' command: %v", err)
	}

	// converts to []string and removes the header
	linesOfProcStrings := strings.Split(strings.TrimSpace(string(output)), "\n")[1:]

	procs := []Proc{}
	for _, line := range linesOfProcStrings {
		pid, err := strconv.Atoi(strings.TrimSpace(line[0:10]))
		if err != nil {
			log.Printf("failed to convert first field to int: %v. split: %v", err, line)
		}
		cpu, err := strconv.ParseFloat(utils.ConvertLocalizedString(strings.TrimSpace(line[63:68])), 64)
		if err != nil {
			log.Printf("failed to convert third field to float: %v. split: %v", err, line)
		}
		mem, err := strconv.ParseFloat(utils.ConvertLocalizedString(strings.TrimSpace(line[69:74])), 64)
		if err != nil {
			log.Printf("failed to convert fourth field to float: %v. split: %v", err, line)
		}
		proc := Proc{
			Pid:         pid,
			CommandName: strings.TrimSpace(line[11:61]),
			Cpu:         cpu,
			Mem:         mem,
			FullCommand: line[74:],
		}
		procs = append(procs, proc)
	}

	return procs, nil
}
