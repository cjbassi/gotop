// +build freebsd

package devices

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/xxxserxxx/gotop/v4/utils"
)

func init() {
	RegisterTemp(update)
	RegisterDeviceList(Temperatures, devs, devs)
}

var sensorOIDS = map[string]string{
	"dev.cpu.0.temperature":           "CPU 0 ",
	"hw.acpi.thermal.tz0.temperature": "Thermal zone 0",
}

func update(temps map[string]int) map[string]error {
	var errors map[string]error

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
	for k, _ := range sensorOIDS {
		rv = append(rv, k)
	}
	return rv
}
