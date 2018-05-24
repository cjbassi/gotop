package widgets

import (
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		self.Data[sensor.SensorKey] = int(sensor.Temperature)
	}
}
