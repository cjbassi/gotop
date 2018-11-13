package widgets

import (
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		if self.Fahrenheit {
			self.Data[sensor.SensorKey] = int(sensor.Temperature*9/5 + 32)
		} else {
			self.Data[sensor.SensorKey] = int(sensor.Temperature)
		}
	}
}
