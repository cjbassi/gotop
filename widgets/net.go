package widgets

import (
	"time"

	ui "github.com/cjbassi/gotop/termui"
	ps "github.com/shirou/gopsutil/net"
)

type Net struct {
	*ui.Sparklines
	interval time.Duration
}

func NewNet() *Net {
	recv := ui.NewSparkline()
	recv.Title = "Receiving"
	recv.Data = []int{0}
	recv.Total = 0

	sent := ui.NewSparkline()
	sent.Title = "Transfering"
	sent.Data = []int{0}
	sent.Total = 0

	spark := ui.NewSparklines(recv, sent)
	n := &Net{spark, time.Second}
	n.Label = "Network Usage"

	go n.update()
	ticker := time.NewTicker(n.interval)
	go func() {
		for range ticker.C {
			n.update()
		}
	}()

	return n
}

func (n *Net) update() {
	interfaces, _ := ps.IOCounters(false)
	recv := int(interfaces[0].BytesRecv)
	sent := int(interfaces[0].BytesSent)

	if n.Lines[0].Total != 0 { // if this isn't the first update
		curRecv := recv - n.Lines[0].Total
		curSent := sent - n.Lines[1].Total

		n.Lines[0].Data = append(n.Lines[0].Data, curRecv)
		n.Lines[1].Data = append(n.Lines[1].Data, curSent)
	}

	n.Lines[0].Total = recv
	n.Lines[1].Total = sent
}
