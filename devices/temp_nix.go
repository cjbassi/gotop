//go:build linux || darwin
// +build linux darwin

package devices

import (
	"fmt"
	"log"

	smart "github.com/anatol/smart.go"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/host"
)

func init() {
	devs() // Populate the sensorMap
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs, defs)
}

func getTemps(temps map[string]int) map[string]error {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		if _, ok := err.(*host.Warnings); ok {
			// ignore warnings
		} else {
			return map[string]error{"gopsutil host": err}
		}
	}
	for _, sensor := range sensors {
		label := sensorMap[sensor.SensorKey]
		if _, ok := temps[label]; ok {
			temps[label] = int(sensor.Temperature)
		}
	}

	block, err := ghw.Block()
	if err != nil {
		log.Print("error getting block device info")
		return nil
	}

	for _, disk := range block.Disks {
		dev, err := smart.Open("/dev/" + disk.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch sm := dev.(type) {
		case *smart.SataDevice:
			data, _ := sm.ReadSMARTData()
			for _, attr := range data.Attrs {
				if attr.Id == 194 {
					temps[disk.Name+"_"+disk.Model] = int(attr.Value)
				}
			}
		case *smart.NVMeDevice:
			data, _ := sm.ReadSMART()
			temps[disk.Name+"_"+disk.Model] = int(data.Temperature)
		default:
		}
	}
	return nil
}

// Optimization to avoid string manipulation every update
var sensorMap map[string]string
