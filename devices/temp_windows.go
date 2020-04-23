// +build windows

package devices

import (
	psHost "github.com/shirou/gopsutil/host"
)

func init() {
	RegisterTemp(update)
	RegisterDeviceList(Temperatures, devs, devs)
}

func update(temps map[string]int) map[string]error {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return map[string]error{"gopsutil": err}
	}
	for _, sensor := range sensors {
		if _, ok := temps[sensor.SensorKey]; ok {
			temps[sensor.SensorKey] = int(sensor.Temperature + 0.5)
		}
	}
	return nil
}

func devs() []string {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return []string{}
	}
	rv := make([]string, 0, len(sensors))
	for _, sensor := range sensors {
		if sensor.Temperature != 0 {
			rv = append(rv, sensor.SensorKey)
		}
	}
	return rv
}
