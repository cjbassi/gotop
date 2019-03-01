package widgets

import (
	"fmt"
	"log"
	"sync"
	"time"

	psNet "github.com/shirou/gopsutil/net"

	ui "github.com/cjbassi/gotop/src/termui"
	"github.com/cjbassi/gotop/src/utils"
)

type NetWidget struct {
	*ui.SparklineGroup
	updateInterval time.Duration

	// used to calculate recent network activity
	totalBytesRecv uint64
	totalBytesSent uint64
}

func NewNetWidget(renderLock *sync.RWMutex) *NetWidget {
	recvSparkline := ui.NewSparkline()
	recvSparkline.Data = []int{}

	sentSparkline := ui.NewSparkline()
	sentSparkline.Data = []int{}

	spark := ui.NewSparklineGroup(recvSparkline, sentSparkline)
	self := &NetWidget{
		SparklineGroup: spark,
		updateInterval: time.Second,
	}
	self.Title = " Network Usage "

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			renderLock.RLock()
			self.update()
			renderLock.RUnlock()
		}
	}()

	return self
}

func (self *NetWidget) update() {
	interfaces, err := psNet.IOCounters(true)
	if err != nil {
		log.Printf("failed to get network activity from gopsutil: %v", err)
		return
	}

	var totalBytesRecv uint64
	var totalBytesSent uint64
	for _, _interface := range interfaces {
		// ignore VPN interface
		if _interface.Name != "tun0" {
			totalBytesRecv += _interface.BytesRecv
			totalBytesSent += _interface.BytesSent
		}
	}

	var recentBytesRecv uint64
	var recentBytesSent uint64

	if self.totalBytesRecv != 0 { // if this isn't the first update
		recentBytesRecv = totalBytesRecv - self.totalBytesRecv
		recentBytesSent = totalBytesSent - self.totalBytesSent

		if int(recentBytesRecv) < 0 {
			log.Printf("error: negative value for recently received network data from gopsutil. recentBytesRecv: %v", recentBytesRecv)
			// recover from error
			recentBytesRecv = 0
		}
		if int(recentBytesSent) < 0 {
			log.Printf("error: negative value for recently sent network data from gopsutil. recentBytesSent: %v", recentBytesSent)
			// recover from error
			recentBytesSent = 0
		}

		self.Lines[0].Data = append(self.Lines[0].Data, int(recentBytesRecv))
		self.Lines[1].Data = append(self.Lines[1].Data, int(recentBytesSent))
	}

	// used in later calls to update
	self.totalBytesRecv = totalBytesRecv
	self.totalBytesSent = totalBytesSent

	// render widget titles
	for i := 0; i < 2; i++ {
		total, label, recent := func() (uint64, string, uint64) {
			if i == 0 {
				return totalBytesRecv, "RX", recentBytesRecv
			}
			return totalBytesSent, "Tx", recentBytesSent
		}()

		recentConverted, unitRecent := utils.ConvertBytes(uint64(recent))
		totalConverted, unitTotal := utils.ConvertBytes(uint64(total))

		self.Lines[i].Title1 = fmt.Sprintf(" Total %s: %5.1f %s", label, totalConverted, unitTotal)
		self.Lines[i].Title2 = fmt.Sprintf(" %s/s: %9.1f %2s/s", label, recentConverted, unitRecent)
	}
}
