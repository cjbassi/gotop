package widgets

import (
	"strings"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	ps "github.com/shirou/gopsutil/host"
)

type Temp struct {
	*ui.List
	interval time.Duration
}

func NewTemp() *Temp {
	t := &Temp{ui.NewList(), time.Second * 5}
	t.Label = "Temperatures"
	t.Threshold = 80 // temp at which color should change to red

	go t.update()
	ticker := time.NewTicker(t.interval)
	go func() {
		for range ticker.C {
			t.update()
		}
	}()

	return t
}

func (t *Temp) update() {
	sensors, _ := ps.SensorsTemperatures()
	temps := []int{}
	labels := []string{}
	for _, temp := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(temp.SensorKey, "input") {
			temps = append(temps, int(temp.Temperature))
			// removes '_input' from the end of the sensor name
			labels = append(labels, temp.SensorKey[:strings.Index(temp.SensorKey, "_input")])
		}
	}
	t.Data = temps
	t.DataLabels = labels
}
