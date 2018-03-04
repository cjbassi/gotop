package widgets

import (
	"time"

	ui "github.com/cjbassi/gotop/termui"
	psMem "github.com/shirou/gopsutil/mem"
)

type Mem struct {
	*ui.LineGraph
	interval time.Duration
}

func NewMem() *Mem {
	m := &Mem{
		LineGraph: ui.NewLineGraph(),
		interval:  time.Second,
	}
	m.Label = "Memory Usage"
	m.Data["Main"] = []float64{0}
	m.Data["Swap"] = []float64{0}

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
	main, _ := psMem.VirtualMemory()
	swap, _ := psMem.SwapMemory()
	m.Data["Main"] = append(m.Data["Main"], main.UsedPercent)
	m.Data["Swap"] = append(m.Data["Swap"], swap.UsedPercent)
}
