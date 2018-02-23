package widgets

import (
	"fmt"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	"github.com/cjbassi/gotop/utils"
	net "github.com/shirou/gopsutil/net"
)

type Net struct {
	*ui.Sparklines
	interval time.Duration
	// used to calculate recent network activity
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
	// `false` causes psutil to group all network activity
	interfaces, _ := net.IOCounters(false)
	recv := interfaces[0].BytesRecv
	sent := interfaces[0].BytesSent

	if n.recvTotal != 0 { // if this isn't the first update
		curRecv := recv - n.recvTotal
		curSent := sent - n.sentTotal

		n.Lines[0].Data = append(n.Lines[0].Data, int(curRecv))
		n.Lines[1].Data = append(n.Lines[1].Data, int(curSent))
	}

	// used for later calls to update
	n.recvTotal = recv
	n.sentTotal = sent

	for i := 0; i < 2; i++ {
		var method string
		var total uint64
		cur := n.Lines[i].Data[len(n.Lines[i].Data)-1]
		totalUnit := "B"
		curUnit := "B"

		if i == 0 {
			total = recv
			method = "Rx"
		} else {
			total = sent
			method = "Tx"
		}

		if cur >= 1000000 {
			cur = int(utils.BytesToMB(uint64(cur)))
			curUnit = "MB"
		} else if cur >= 1000 {
			cur = int(utils.BytesToKB(uint64(cur)))
			curUnit = "kB"
		}

		var totalCvrt float64
		if total >= 1000000000 {
			totalCvrt = utils.BytesToGB(total)
			totalUnit = "GB"
		} else if total >= 1000000 {
			totalCvrt = utils.BytesToMB(total)
			totalUnit = "MB"
		}

		n.Lines[i].Title1 = fmt.Sprintf(" Total %s: %5.1f %s", method, totalCvrt, totalUnit)
		n.Lines[i].Title2 = fmt.Sprintf(" %s/s: %9d %2s/s", method, cur, curUnit)
	}
}
