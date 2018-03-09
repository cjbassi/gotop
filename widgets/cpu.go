package widgets

import (
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	count    int // number of CPUs
	interval time.Duration
}

func NewCPU(interval time.Duration) *CPU {
	count, _ := psCPU.Counts(false)
	c := &CPU{
		LineGraph: ui.NewLineGraph(),
		count:     count,
		interval:  interval,
	}
	c.Label = "CPU Usage"
	for i := 0; i < c.count; i++ {
		key := "CPU" + strconv.Itoa(i+1)
		c.Data[key] = []float64{0}
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
	// `true` makes it so psutil doesn't group CPU usage percentages
	percent, _ := psCPU.Percent(c.interval, true)
	for i := 0; i < c.count; i++ {
		key := "CPU" + strconv.Itoa(i+1)
		c.Data[key] = append(c.Data[key], percent[i])
	}
}
