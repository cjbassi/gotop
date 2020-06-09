// +build linux

package devices

import (
	"strings"

	"github.com/shirou/gopsutil/host"
)

// All possible thermometers
func devs() []string {
	if sensorMap == nil {
		sensorMap = make(map[string]string)
	}
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return []string{}
	}
	rv := make([]string, 0, len(sensors))
	for _, sensor := range sensors {
		label := sensor.SensorKey
		if strings.Contains(sensor.SensorKey, "input") {
			label = sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
		}
		rv = append(rv, label)
		sensorMap[sensor.SensorKey] = label
	}
	return rv
}

// Only include sensors with input in their name; these are the only sensors
// returning live data
func defs() []string {
	// MUST be called AFTER init()
	rv := make([]string, 0)
	for k, v := range sensorMap {
		if k != v { // then it's an _input sensor
			rv = append(rv, v)
		}
	}
	return rv
}
