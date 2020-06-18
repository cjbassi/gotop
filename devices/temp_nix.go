// +build linux darwin

package devices

import (
	"github.com/shirou/gopsutil/host"
)

func init() {
	devs() // Populate the sensorMap
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs, defs)
}

func getTemps(temps map[string]int) map[string]error {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return map[string]error{"gopsutil host": err}
	}
	for _, sensor := range sensors {
		label := sensorMap[sensor.SensorKey]
		if _, ok := temps[label]; ok {
			temps[label] = int(sensor.Temperature)
		}
	}
	return nil
}

// Optimization to avoid string manipulation every update
var sensorMap map[string]string
