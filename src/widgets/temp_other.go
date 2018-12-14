// +build linux freebsd

package widgets

import (
	"log"
	"strings"

	"github.com/cjbassi/gotop/src/utils"
	psHost "github.com/shirou/gopsutil/host"
)

func (self *Temp) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		log.Printf("failed to get sensors from gopsutil: %v", err)
	}
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			if self.Fahrenheit {
				self.Data[label] = utils.CelsiusToFahrenheit(int(sensor.Temperature))
			} else {
				self.Data[label] = int(sensor.Temperature)
			}
		}
	}
}
