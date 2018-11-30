package widgets

import (
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		if sensor.Temperature != 0 {
			self.Data[sensor.SensorKey] = int(sensor.Temperature)
		}
	}
}
