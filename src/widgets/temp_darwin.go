package widgets

// #cgo LDFLAGS: -framework IOKit
// #include "include/smc.c"
import "C"
import (
	"log"

	"github.com/cjbassi/gotop/src/utils"
)

type TemperatureStat struct {
	SensorKey   string  `json:"sensorKey"`
	Temperature float64 `json:"sensorTemperature"`
}

func SensorsTemperatures() ([]TemperatureStat, error) {
	temperatureKeys := map[string]string{
		C.AMBIENT_AIR_0:          "ambient_air_0",
		C.AMBIENT_AIR_1:          "ambient_air_1",
		C.CPU_0_DIODE:            "cpu_0_diode",
		C.CPU_0_HEATSINK:         "cpu_0_heatsink",
		C.CPU_0_PROXIMITY:        "cpu_0_proximity",
		C.ENCLOSURE_BASE_0:       "enclosure_base_0",
		C.ENCLOSURE_BASE_1:       "enclosure_base_1",
		C.ENCLOSURE_BASE_2:       "enclosure_base_2",
		C.ENCLOSURE_BASE_3:       "enclosure_base_3",
		C.GPU_0_DIODE:            "gpu_0_diode",
		C.GPU_0_HEATSINK:         "gpu_0_heatsink",
		C.GPU_0_PROXIMITY:        "gpu_0_proximity",
		C.HARD_DRIVE_BAY:         "hard_drive_bay",
		C.MEMORY_SLOT_0:          "memory_slot_0",
		C.MEMORY_SLOTS_PROXIMITY: "memory_slots_proximity",
		C.NORTHBRIDGE:            "northbridge",
		C.NORTHBRIDGE_DIODE:      "northbridge_diode",
		C.NORTHBRIDGE_PROXIMITY:  "northbridge_proximity",
		C.THUNDERBOLT_0:          "thunderbolt_0",
		C.THUNDERBOLT_1:          "thunderbolt_1",
		C.WIRELESS_MODULE:        "wireless_module",
	}

	var temperatures []TemperatureStat

	C.open_smc()
	defer C.close_smc()

	for key, val := range temperatureKeys {
		temperatures = append(temperatures, TemperatureStat{
			SensorKey:   val,
			Temperature: float64(C.get_tmp(C.CString(key), C.CELSIUS)),
		})
	}
	return temperatures, nil
}

func (self *TempWidget) update() {
	sensors, err := SensorsTemperatures()
	if err != nil {
		log.Printf("failed to get sensors from CGO: %v", err)
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
