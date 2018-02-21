package widgets

import (
	"time"

	ui "github.com/cjbassi/gotop/termui"
	ps "github.com/shirou/gopsutil/mem"
)

type Mem struct {
	*ui.LineGraph
	interval time.Duration
}

func NewMem() *Mem {
	m := &Mem{ui.NewLineGraph(), time.Second}
	m.Label = "Memory Usage"
	m.Data["Main"] = []float64{0} // Sets initial data to 0
	m.Data["Swap"] = []float64{0}
	m.LineColor["Main"] = ui.Color(5)
	m.LineColor["Swap"] = ui.Color(11)

	go m.update()
	ticker := time.NewTicker(m.interval)
	go func() {
		for range ticker.C {
			m.update()
		}
	}()

	return m
}

func (m *Mem) update() {
	main, _ := ps.VirtualMemory()
	swap, _ := ps.SwapMemory()
	m.Data["Main"] = append(m.Data["Main"], main.UsedPercent)
	m.Data["Swap"] = append(m.Data["Swap"], swap.UsedPercent)
}
