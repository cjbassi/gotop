package widgets

import (
	"fmt"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psNet "github.com/shirou/gopsutil/net"
)

type Net struct {
	*ui.Sparklines
	interval time.Duration
	// used to calculate recent network activity
	prevRecvTotal uint64
	prevSentTotal uint64
}

func NewNet() *Net {
	recv := ui.NewSparkline()
	recv.Data = []int{0}

	sent := ui.NewSparkline()
	sent.Data = []int{0}

	spark := ui.NewSparklines(recv, sent)
	self := &Net{
		Sparklines: spark,
		interval:   time.Second,
	}
	self.Label = "Network Usage"

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Net) update() {
	// `false` causes psutil to group all network activity
	interfaces, _ := psNet.IOCounters(false)
	curRecvTotal := interfaces[0].BytesRecv
	curSentTotal := interfaces[0].BytesSent

	if self.prevRecvTotal != 0 { // if this isn't the first update
		recvRecent := curRecvTotal - self.prevRecvTotal
		sentRecent := curSentTotal - self.prevSentTotal

		self.Lines[0].Data = append(self.Lines[0].Data, int(recvRecent))
		self.Lines[1].Data = append(self.Lines[1].Data, int(sentRecent))

		if int(recvRecent) < 0 || int(sentRecent) < 0 {
			utils.Error("net data",
				fmt.Sprint(
					"curRecvTotal: ", curRecvTotal, "\n",
					"curSentTotal: ", curSentTotal, "\n",
					"self.prevRecvTotal: ", self.prevRecvTotal, "\n",
					"self.prevSentTotal: ", self.prevSentTotal, "\n",
					"recvRecent: ", recvRecent, "\n",
					"sentRecent: ", sentRecent, "\n",
					"int(recvRecent): ", int(recvRecent), "\n",
					"int(sentRecent): ", int(sentRecent),
				))
		}
	}

	// used in later calls to update
	self.prevRecvTotal = curRecvTotal
	self.prevSentTotal = curSentTotal

	// render widget titles
	for i := 0; i < 2; i++ {
		var method string // either 'Rx' or 'Tx'
		var total float64
		recent := self.Lines[i].Data[len(self.Lines[i].Data)-1]

		if i == 0 {
			total = float64(curRecvTotal)
			method = "Rx"
		} else {
			total = float64(curSentTotal)
			method = "Tx"
		}

		recentFloat, unitRecent := utils.ConvertBytes(uint64(recent))
		recent = int(recentFloat)
		total, unitTotal := utils.ConvertBytes(uint64(total))

		self.Lines[i].Title1 = fmt.Sprintf(" Total %s: %5.1f %s", method, total, unitTotal)
		self.Lines[i].Title2 = fmt.Sprintf(" %s/s: %9d %2s/s", method, recent, unitRecent)
	}
}
