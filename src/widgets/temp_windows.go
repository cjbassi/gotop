package widgets

import (
	"log"

	psHost "github.com/shirou/gopsutil/host"

	"github.com/cjbassi/gotop/src/utils"
)

func (self *TempWidget) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		log.Printf("failed to get sensors from gopsutil: %v", err)
		return
	}
	for _, sensor := range sensors {
		if sensor.Temperature != 0 {
			switch self.TempScale {
			case Fahrenheit:
				self.Data[sensor.SensorKey] = utils.CelsiusToFahrenheit(int(sensor.Temperature))
			case Celcius:
				self.Data[sensor.SensorKey] = int(sensor.Temperature)
			}
		}
	}
}
