// +build darwin

package devices

import (
	"bytes"
	"encoding/csv"
	"github.com/shirou/gopsutil/host"
	"io"
)

// All possible thermometers
func devs() []string {
	// Did we already populate the sensorMap?
	if sensorMap != nil {
		return defs()
	}
	// Otherwise, get the sensor data from the system & filter it
	ids := loadIDs()
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		// FIXME log an error here
		return []string{}
	}
	rv := make([]string, 0, len(sensors))
	sensorMap = make(map[string]string)
	for _, sensor := range sensors {
		// 0-value sensors are not implemented
		if sensor.Temperature == 0 {
			continue
		}
		if label, ok := ids[sensor.SensorKey]; ok {
    			sensorMap[sensor.SensorKey] = label
			rv = append(rv, label)
		}
	}
	return rv
}

// Only the ones filtered
func defs() []string {
	rv := make([]string, 0, len(sensorMap))
	for _, val := range sensorMap {
		rv = append(rv, val)
	}
	return rv
}

// loadIDs parses the embedded smc.tsv data that maps Darwin SMC
// sensor IDs to their human-readable labels into an array and returns the
// array. The array keys are the 4-letter sensor keys; the values are the
// human labels.
func loadIDs() map[string]string {
	rv := make(map[string]string)
	data, err := Asset("smc.tsv")
	parser := csv.NewReader(bytes.NewReader(data))
	parser.Comma = '\t'
	var line []string
	for {
		if line, err = parser.Read(); err == io.EOF {
			break
		}
		if err != nil {
			// FIXME log an error here
			break
		}
		// The line is malformed if len(line) != 2, but because the asset is static
		// it makes no sense to report the error to downstream users. This must be
		// tested at/around compile time.
		// FIXME assert all lines in smc.tsv have 2 columns during unit tests
		if len(line) == 2 {
			rv[line[0]] = line[1]
		}
	}
	return rv
}
