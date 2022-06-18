//go:build linux || darwin
// +build linux darwin

package devices

import (
	"log"

	"github.com/anatol/smart.go"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/block"
	"github.com/shirou/gopsutil/host"
)

var disks []block.Disk
var smDevices []smart.Device

func init() {
	devs() // Populate the sensorMap
	RegisterStartup(startBlock)
	RegisterTemp(getTemps)
	RegisterDeviceList(Temperatures, devs, defs)
	RegisterShutdown(endBlock)
}

func startBlock(vars map[string]string) error {
	block, err := ghw.Block()
	if err != nil {
		log.Print("error getting block device info")
		return err
	}
	for _, disk := range block.Disks {
		dev, err := smart.Open("/dev/" + disk.Name)
		if err == nil {
			disks = append(disks, *disk)
			smDevices = append(smDevices, dev)
		}
	}
	return err
}

func endBlock() error {
	for _, dev := range smDevices {
		err := dev.Close()
		if err != nil {
			log.Print("error closing device")
			return err
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

	for i, dev := range smDevices {
		switch sm := dev.(type) {
		case *smart.SataDevice:
			data, err := sm.ReadSMARTData()
			if err != nil {
				log.Print("error getting smart data for " + disks[i].Name + "_" + disks[i].Model)
				log.Print(err)
				break
			}
			if attr, ok := data.Attrs[194]; ok {
				val, _, _, _, err := attr.ParseAsTemperature()
				if err != nil {
					log.Print("error parsing temperature smart data for " + disks[i].Name + "_" + disks[i].Model)
					log.Print(err)
					break
				}
				temps[disks[i].Name+"_"+disks[i].Model] = int(val)
			}
		case *smart.NVMeDevice:
			data, err := sm.ReadSMART()
			if err != nil {
				log.Print("error getting smart data for " + disks[i].Name + "_" + disks[i].Model)
				break
			}
			temps[disks[i].Name+"_"+disks[i].Model] = int(data.Temperature)
		default:
		}
	}
	return nil
}

// Optimization to avoid string manipulation every update
var sensorMap map[string]string
