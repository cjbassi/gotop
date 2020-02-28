// +build freebsd

package devices

import (
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	mf := func(mems map[string]MemoryInfo) map[string]error {
		cmd := "swapinfo -k|sed -n '1!p'|awk '{print $2,$3,$5}'"
		output, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			return map[string]error{"swapinfo": err}
		}

		s := strings.TrimSuffix(string(output), "\n")
		s = strings.ReplaceAll(s, "\n", " ")
		ss := strings.Split(s, " ")
		ss = ss[((len(ss)/3)-1)*3:]

		errors := make(map[string]error)
		mem := MemoryInfo{}
		mem.Total, err = strconv.ParseUint(ss[0], 10, 64)
		if err != nil {
			errors["swap total"] = err
		}

		mem.Used, err = strconv.ParseUint(ss[1], 10, 64)
		if err != nil {
			errors["swap used"] = err
		}

		mem.UsedPercent, err = strconv.ParseFloat(strings.TrimSuffix(ss[2], "%"), 64)
		if err != nil {
			errors["swap percent"] = err
		}
		mems["Swap"] = mem
		return errors
	}
	RegisterMem(mf)
}
