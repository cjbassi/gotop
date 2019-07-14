package widgets

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
	"github.com/rai-project/nvidia-smi"
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

func collectSysctlSensors() []sensorMeasurement {
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

	return measurements
}

func collectNvidiaSensors() []sensorMeasurement {
	var measurements []sensorMeasurement

	info, _ := nvidiasmi.New()
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			var s sensorMeasurement
			s.name = strings.ReplaceAll(strings.ToLower(gpu.ProductName), " ", "_") + "_" + strconv.Itoa(i) + "_input"
			s.temperature, _ = strconv.ParseFloat(strings.ReplaceAll(gpu.GpuTemp, " C", ""), 10)
			measurements = append(measurements, s)
		}
	}

	return measurements
}

func collectAMDGPUSensors() []sensorMeasurement {
	var measurments []sensorMeasurement

	return measurments
}

func collectGPUSensors() []sensorMeasurement {
	var measurements []sensorMeasurement

	measurements = append(measurements, collectSysctlSensors()...)
	measurements = append(measurements, collectNvidiaSensors()...)
	measurements = append(measurements, collectAMDGPUSensors()...)

	return measurements
}

func collectSensors() []sensorMeasurement {
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

	measurements = append(measurements, collectGPUSensors()...)

	return measurements

}

func (self *TempWidget) update() {
	sensors := collectSensors()

	for _, sensor := range sensors {
		switch self.TempScale {
		case Fahrenheit:
			self.Data[sensor.name] = utils.CelsiusToFahrenheit(int(sensor.temperature))
		case Celcius:
			self.Data[sensor.name] = int(sensor.temperature)
		}
	}
}
