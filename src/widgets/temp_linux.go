package widgets

import (
	"strings"

	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			if self.Fahrenheit {
				self.Data[label] = int(sensor.Temperature*9/5 + 32)
			} else {
				self.Data[label] = int(sensor.Temperature)
			}
		}
	}
}
