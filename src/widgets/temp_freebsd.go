// +build freebsd

package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

var sensorOIDS = []string{
	"dev.cpu.0.temperature",
	"hw.acpi.thermal.tz0.temperature",
}

type sensor struct {
	name        string
	temperature float64
}

func removeUnusedChars(s string) string {
	s1 := strings.Replace(s, "C", "", 1)
	s2 := strings.TrimSuffix(s1, "\n")
	return s2
}

func refineOutput(output []byte) (float64, error) {
	convertedOutput := utils.ConvertLocalizedString(removeUnusedChars(string(output)))
	value, err := strconv.ParseFloat(convertedOutput, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func collectSensors() ([]sensor, error) {
	sensorData := []sensor{}
	for _, v := range sensorOIDS {
		output, err := exec.Command("sysctl", "-n", v).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'sysctl' command: %v", err)
		}

		value, err := refineOutput(output)
		if err != nil {
			return nil, fmt.Errorf("failed to execute 'sysctl' command: %v", err)
		}

		sensorData = append(sensorData, sensor{v, value})

	}
	return sensorData, nil

}

func (self *TempWidget) update() {
	sensors, err := collectSensors()
	if err != nil {
		log.Printf("error recieved from gopsutil: %v", err)
	}
	for _, sensor := range sensors {
		switch self.TempScale {
		case Fahrenheit:
			self.Data[sensor.name] = utils.CelsiusToFahrenheit(int(sensor.temperature))
		case Celcius:
			self.Data[sensor.name] = int(sensor.temperature)
		}
	}
}
