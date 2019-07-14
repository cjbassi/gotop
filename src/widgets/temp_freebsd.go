package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
)

var sensorOIDS = map[string]string{
	"dev.cpu.0.temperature":           "CPU 0 ",
	"hw.acpi.thermal.tz0.temperature": "Thermal zone 0",
}

type sensorMeasurement struct {
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

func collectNvidiaGPUSensors() []sensorMeasurement {
	var measurements []sensorMeasurement

	_, err := exec.Command("sysctl", "-n", "hw.nvidia.gpus.0.model").Output()
	if err == nil {
		output, err := exec.Command("nvidia-settings", "-q", "gpucoretemp", "-t").Output()
		if err != nil {
			log.Printf("Failed to get nvidia gpu temperature from nvidia-settings")
		} else {
			s := strings.TrimSuffix(string(output), "\n")
			s = strings.ReplaceAll(s, "\n", " ")
			ss := strings.Split(s, " ")
			for i := 1; i < len(ss); i++ {
				value, err := refineOutput([]byte(ss[i]))
				if err != nil {
					log.Printf("Failed to parse nvidia-settings output")
				}
				label := fmt.Sprintf("Nvidia GPU %d", i-1)
				measurements = append(measurements, sensorMeasurement{label, value})
			}
		}
	}

	return measurements
}

func collectGPUSensors() []sensorMeasurement {
	var measurements []sensorMeasurement

	m := collectNvidiaGPUSensors()
	measurements = append(measurements, m...)

	return measurements
}

func collectSensors() ([]sensorMeasurement, error) {
	var measurements []sensorMeasurement
	for k, v := range sensorOIDS {
		output, err := exec.Command("sysctl", "-n", k).Output()
		if err != nil {
			continue
		}

		value, err := refineOutput(output)
		if err != nil {
			continue
		}

		measurements = append(measurements, sensorMeasurement{v, value})
	}

	m := collectGPUSensors()
	measurements = append(measurements, m...)

	return measurements, nil

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
