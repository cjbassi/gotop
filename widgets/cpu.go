package widgets

import (
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	ps "github.com/shirou/gopsutil/cpu"
)

type CPU struct {
	*ui.LineGraph
	count    int
	interval time.Duration
}

func NewCPU() *CPU {
	count, _ := ps.Counts(false)
	c := &CPU{ui.NewLineGraph(), count, time.Second}
	c.Label = "CPU Usage"
	for i := 0; i < c.count; i++ {
		key := "CPU" + strconv.Itoa(i+1)
		c.Data[key] = []float64{0}
		c.LineColor[key] = ui.Attribute(int(ui.ColorRed) + i)
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
	percent, _ := ps.Percent(time.Second, true) // takes one second to get the data
	for i := 0; i < c.count; i++ {
		key := "CPU" + strconv.Itoa(i+1)
		c.Data[key] = append(c.Data[key], percent[i])
	}
}
