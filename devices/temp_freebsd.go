// +build freebsd

package devices

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/xxxserxxx/gotop/v4/utils"
)

func init() {
	if len(devs()) == 0 {
		log.Println(tr.Value("error.nodevfound", "thermal sensors"))
		return
	}
	RegisterTemp(update)
	RegisterDeviceList(Temperatures, devs, devs)
}

var sensorOIDS = map[string]string{
	"dev.cpu.0.temperature":           "CPU 0 ",
	"hw.acpi.thermal.tz0.temperature": "Thermal zone 0",
}

func update(temps map[string]int) map[string]error {
	errors := make(map[string]error)

	for k, v := range sensorOIDS {
		if _, ok := temps[k]; !ok {
			continue
		}
		output, err := exec.Command("sysctl", "-n", k).Output()
		if err != nil {
			errors[v] = err
			continue
		}

		s1 := strings.Replace(string(output), "C", "", 1)
		s2 := strings.TrimSuffix(s1, "\n")
		convertedOutput := utils.ConvertLocalizedString(s2)
		value, err := strconv.ParseFloat(convertedOutput, 64)
		if err != nil {
			errors[v] = err
			continue
		}

		temps[v] = int(value)
	}

	return errors
}

func devs() []string {
	rv := make([]string, 0, len(sensorOIDS))
	// Check that thermal sensors are really available; they aren't in VMs
	bs, err := exec.Command("sysctl", "-a").Output()
	if err != nil {
		log.Printf(tr.Value("error.fatalfetch", "temp", err.Error()))
		return []string{}
	}
	for k, _ := range sensorOIDS {
		idx := strings.Index(string(bs), k)
		if idx >= 0 {
			rv = append(rv, k)
		}
	}
	if len(rv) == 0 {
		oids := make([]string, 0, len(sensorOIDS))
		for k, _ := range sensorOIDS {
			oids = append(oids, k)
		}
		log.Printf(tr.Value("error.nodevfound", strings.Join(oids, ", ")))
	}
	return rv
}
