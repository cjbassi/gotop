// +build darwin

package devices

// #cgo LDFLAGS: -framework IOKit
// #include "include/smc.c"
import "C"

func init() {
	RegisterTemp(update)
}

func update(temps map[string]int) map[string]error {
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

	C.open_smc()
	defer C.close_smc()

	for key, val := range temperatureKeys {
		temps[val] = int(C.get_tmp(C.CString(key), C.CELSIUS))
	}

	return nil
}
