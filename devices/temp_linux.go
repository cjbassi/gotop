// +build linux

package devices

import (
	"strings"

	psHost "github.com/shirou/gopsutil/host"
)

func init() {
	devs() // Populate the sensorMap
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs, defs)
}

func getTemps(temps map[string]int) map[string]error {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return map[string]error{"psHost": err}
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

func devs() []string {
	if sensorMap == nil {
		sensorMap = make(map[string]string)
	}
	sensors, err := psHost.SensorsTemperatures()
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
