// +build freebsd

package devices

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/xxxserxxx/gotop/v3/utils"
)

func init() {
	RegisterTemp(update)
}

var sensorOIDS = map[string]string{
	"dev.cpu.0.temperature":           "CPU 0 ",
	"hw.acpi.thermal.tz0.temperature": "Thermal zone 0",
}

func update(temps map[string]int) map[string]error {
	errors := make(map[string]error)

	for k, v := range sensorOIDS {
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
