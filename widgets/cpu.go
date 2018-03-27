package widgets

import (
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	Count    int // number of CPUs
	interval time.Duration
}

func NewCPU(interval time.Duration, zoom int) *CPU {
	count, _ := psCPU.Counts(false)
	c := &CPU{
		LineGraph: ui.NewLineGraph(),
		Count:     count,
		interval:  interval,
	}
	c.Label = "CPU Usage"
	c.Zoom = zoom
	if c.Count <= 8 {
		for i := 0; i < c.Count; i++ {
			key := "CPU" + strconv.Itoa(i+1)
			c.Data[key] = []float64{0}
		}
	} else {
		c.Data["Average"] = []float64{0}
	}

	go c.update()
	ticker := time.NewTicker(c.interval)
	go func() {
		for range ticker.C {
			c.update()
		}
	}()

	return c
}

func (c *CPU) update() {
	// psutil calculates the CPU usage over a 1 second interval, therefore it blocks for 1 second
	if c.Count <= 8 {
		percent, _ := psCPU.Percent(c.interval, true)
		for i := 0; i < c.Count; i++ {
			key := "CPU" + strconv.Itoa(i+1)
			c.Data[key] = append(c.Data[key], percent[i])
		}
	} else {
		percent, _ := psCPU.Percent(c.interval, false)
		c.Data["Average"] = append(c.Data["Average"], percent[0])
	}
}
