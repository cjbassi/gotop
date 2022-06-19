//go:build linux || darwin
// +build linux darwin

package devices

import (
	"log"

	"github.com/anatol/smart.go"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/host"
)

var smDevices map[string]smart.Device

func init() {
	devs() // Populate the sensorMap
	RegisterStartup(startBlock)
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs, defs)
	RegisterShutdown(endBlock)
}

func startBlock(vars map[string]string) error {
	smDevices = make(map[string]smart.Device)

	block, err := ghw.Block()
	if err != nil {
		log.Printf("error getting block device info: %s", err)
		return err
	}
	for _, disk := range block.Disks {
		dev, err := smart.Open("/dev/" + disk.Name)
		if err != nil {
			log.Printf("error opening smart info for %s: %s", disk.Name, err)
			continue
		}
		smDevices[disk.Name+"_"+disk.Model] = dev
	}
	return nil
}

func endBlock() error {
	for name, dev := range smDevices {
		err := dev.Close()
		if err != nil {
			log.Printf("error closing device %s: %s", name, err)
		}
	}
	return nil
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

	for name, dev := range smDevices {
		switch sm := dev.(type) {
		case *smart.SataDevice:
			data, err := sm.ReadSMARTData()
			if err != nil {
				log.Printf("error getting smart data for %s: %s", name, err)
				continue
			}
			if attr, ok := data.Attrs[194]; ok {
				val, _, _, _, err := attr.ParseAsTemperature()
				if err != nil {
					log.Printf("error parsing temperature smart data for %s: %s", name, err)
					continue
				}
				temps[name] = val
			}
		case *smart.NVMeDevice:
			data, err := sm.ReadSMART()
			if err != nil {
				log.Printf("error getting smart data for %s: %s", name, err)
				continue
			}
			temps[name] = int(data.Temperature)
		default:
		}
	}
	return nil
}

// Optimization to avoid string manipulation every update
var sensorMap map[string]string
