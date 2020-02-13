package widgets

import (
	"log"
	"strings"

	psHost "github.com/shirou/gopsutil/host"

	"github.com/cjbassi/gotop/utils"
)

func (self *TempWidget) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		log.Printf("error received from gopsutil: %v", err)
	}
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			switch self.TempScale {
			case Fahrenheit:
				self.Data[label] = utils.CelsiusToFahrenheit(int(sensor.Temperature))
			case Celsius:
				self.Data[label] = int(sensor.Temperature)
			}
		}
	}
}
