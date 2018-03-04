package widgets

import (
	"fmt"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	"github.com/cjbassi/gotop/utils"
	psNet "github.com/shirou/gopsutil/net"
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
	n := &Net{
		Sparklines: spark,
		interval:   time.Second,
	}
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
	interfaces, _ := psNet.IOCounters(false)
	recvTotal := interfaces[0].BytesRecv
	sentTotal := interfaces[0].BytesSent

	if n.recvTotal != 0 { // if this isn't the first update
		recvRecent := recvTotal - n.recvTotal
		sentRecent := sentTotal - n.sentTotal

		n.Lines[0].Data = append(n.Lines[0].Data, int(recvRecent))
		n.Lines[1].Data = append(n.Lines[1].Data, int(sentRecent))
	}

	// used in later calls to update
	n.recvTotal = recvTotal
	n.sentTotal = sentTotal

	// renders net widget titles
	for i := 0; i < 2; i++ {
		var method string // either 'Rx' or 'Tx'
		var total float64
		recent := n.Lines[i].Data[len(n.Lines[i].Data)-1]
		unitTotal := "B"
		unitRecent := "B"

		if i == 0 {
			total = float64(recvTotal)
			method = "Rx"
		} else {
			total = float64(sentTotal)
			method = "Tx"
		}

		if recent >= 1000000 {
			recent = int(utils.BytesToMB(uint64(recent)))
			unitRecent = "MB"
		} else if recent >= 1000 {
			recent = int(utils.BytesToKB(uint64(recent)))
			unitRecent = "kB"
		}

		if total >= 1000000000 {
			total = utils.BytesToGB(uint64(total))
			unitTotal = "GB"
		} else if total >= 1000000 {
			total = utils.BytesToMB(uint64(total))
			unitTotal = "MB"
		}

		n.Lines[i].Title1 = fmt.Sprintf(" Total %s: %5.1f %s", method, total, unitTotal)
		n.Lines[i].Title2 = fmt.Sprintf(" %s/s: %9d %2s/s", method, recent, unitRecent)
	}
}
