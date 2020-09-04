package widgets

import (
	"log"
	"strings"

	psHost "github.com/shirou/gopsutil/host"

	"github.com/cjbassi/gotop/src/utils"
	"github.com/rai-project/nvidia-smi"
	"strconv"
)

func (self *TempWidget) update() {
	sensors, err := psHost.SensorsTemperatures()
	if err != nil {
		log.Printf("error recieved from gopsutil: %v", err)
	}
	info, err := nvidiasmi.New()

	if err != nil {
		log.Printf("Error while requesting GPU: %s", err.Error())
	}
	if info.HasGPU() {
		for i := range info.GPUS {
			gpu := info.GPUS[i]
			var gt psHost.TemperatureStat
			gt.SensorKey = strings.ReplaceAll(strings.ToLower(gpu.ProductName), " ", "_") + "_" + strconv.Itoa(i) + "_input"
			gt.Temperature, _ = strconv.ParseFloat(strings.ReplaceAll(gpu.GpuTemp, " C", ""), 10)
			sensors = append(sensors, gt)
		}
	}

	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") && sensor.Temperature != 0 {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			switch self.TempScale {
			case Fahrenheit:
				self.Data[label] = utils.CelsiusToFahrenheit(int(sensor.Temperature))
			case Celcius:
				self.Data[label] = int(sensor.Temperature)
			}
		}
	}
}
