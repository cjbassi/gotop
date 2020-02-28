// +build windows

package devices

import (
	psHost "github.com/shirou/gopsutil/host"
)

func init() {
	RegisterTemp(update)
}

func update(temps map[string]int) map[string]error {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		return map[string]error{"gopsutil": err}
	}
	for _, sensor := range sensors {
		if sensor.Temperature != 0 {
			temps[sensor.SensorKey] = int(sensor.Temperature)
		}
	}
	return nil
}
