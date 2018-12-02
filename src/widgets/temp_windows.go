package widgets

import (
	"github.com/cjbassi/gotop/src/utils"
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		if sensor.Temperature != 0 {
			if self.Fahrenheit {
				self.Data[sensor.SensorKey] = utils.CelsiusToFahrenheit(int(sensor.Temperature))
			} else {
				self.Data[sensor.SensorKey] = int(sensor.Temperature)
			}
		}
	}
}
