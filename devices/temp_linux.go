// +build linux

package devices

import (
	"strings"

	psHost "github.com/shirou/gopsutil/host"
)

func init() {
	RegisterTemp(getTemps)
}

func getTemps(temps map[string]int) map[string]error {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return map[string]error{"psHost": err}
	}
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			temps[label] = int(sensor.Temperature)
		}
	}
	return nil
}
