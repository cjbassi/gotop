package widgets

import (
	"fmt"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	"github.com/cjbassi/gotop/utils"
	ps "github.com/shirou/gopsutil/net"
)

type Net struct {
	*ui.Sparklines
	interval  time.Duration
	recvTotal uint64
	sentTotal uint64
}

func NewNet() *Net {
	recv := ui.NewSparkline()
	recv.Data = []int{0}

	sent := ui.NewSparkline()
	sent.Data = []int{0}

	spark := ui.NewSparklines(recv, sent)
	n := &Net{spark, time.Second, 0, 0}
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
	recv := interfaces[0].BytesRecv
	sent := interfaces[0].BytesSent

	if n.recvTotal != 0 { // if this isn't the first update
		curRecv := recv - n.recvTotal
		curSent := sent - n.sentTotal

		n.Lines[0].Data = append(n.Lines[0].Data, int(curRecv))
		n.Lines[1].Data = append(n.Lines[1].Data, int(curSent))
	}

	n.recvTotal = recv
	n.sentTotal = sent

	for i := 0; i < 2; i++ {
		var method string
		var total uint64
		cur := n.Lines[i].Data[len(n.Lines[i].Data)-1]
		totalUnit := "B"
		curUnit := "B"

		if i == 0 {
			total = n.recvTotal
			method = "Rx"
		} else {
			total = n.sentTotal
			method = "Tx"
		}

		if cur >= 1000000 {
			cur = int(utils.BytesToMB(uint64(cur)))
			curUnit = "MB"
		} else if cur >= 1000 {
			cur = int(utils.BytesToKB(uint64(cur)))
			curUnit = "kB"
		}

		if total >= 1000000000 {
			total = utils.BytesToGB(total)
			totalUnit = "GB"
		} else if total >= 1000000 {
			total = utils.BytesToMB(total)
			totalUnit = "MB"
		}

		n.Lines[i].Title1 = fmt.Sprintf(" Total %s: %3d %s", method, total, totalUnit)
		n.Lines[i].Title2 = fmt.Sprintf(" %s/s: %7d %2s/s", method, cur, curUnit)
	}
}
