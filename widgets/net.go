package widgets

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	psNet "github.com/shirou/gopsutil/net"

	ui "github.com/xxxserxxx/gotop/v4/termui"
	"github.com/xxxserxxx/gotop/v4/utils"
)

const (
	// NetInterfaceAll enables all network interfaces
	NetInterfaceAll = "all"
	// NetInterfaceVpn is the VPN interface
	NetInterfaceVpn = "tun0"
)

type NetWidget struct {
	*ui.SparklineGroup
	updateInterval time.Duration

	// used to calculate recent network activity
	totalBytesRecv uint64
	totalBytesSent uint64
	NetInterface   []string
	sentMetric     prometheus.Counter
	recvMetric     prometheus.Counter
	Mbps           bool
}

// TODO: state:merge #169 % option for network use (jrswab/networkPercentage)
func NewNetWidget(netInterface string) *NetWidget {
	recvSparkline := ui.NewSparkline()
	recvSparkline.Data = []int{}

	sentSparkline := ui.NewSparkline()
	sentSparkline.Data = []int{}

	spark := ui.NewSparklineGroup(recvSparkline, sentSparkline)
	self := &NetWidget{
		SparklineGroup: spark,
		updateInterval: time.Second,
		NetInterface:   strings.Split(netInterface, ","),
	}
	self.Title = " Network Usage "
	if netInterface != "all" {
		self.Title = fmt.Sprintf(" Network Usage: %s ", netInterface)
	}

	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (net *NetWidget) EnableMetric() {
	net.recvMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gotop",
		Subsystem: "net",
		Name:      "recv",
	})
	prometheus.MustRegister(net.recvMetric)

	net.sentMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "gotop",
		Subsystem: "net",
		Name:      "sent",
	})
	prometheus.MustRegister(net.sentMetric)
}

func (net *NetWidget) update() {
	interfaces, err := psNet.IOCounters(true)
	if err != nil {
		log.Printf("failed to get network activity from gopsutil: %v", err)
		return
	}

	var totalBytesRecv uint64
	var totalBytesSent uint64
	interfaceMap := make(map[string]bool)
	// Default behaviour
	interfaceMap[NetInterfaceAll] = true
	interfaceMap[NetInterfaceVpn] = false
	// Build a map with wanted status for each interfaces.
	for _, iface := range net.NetInterface {
		if strings.HasPrefix(iface, "!") {
			interfaceMap[strings.TrimPrefix(iface, "!")] = false
		} else {
			// if we specify a wanted interface, remove capture on all.
			delete(interfaceMap, NetInterfaceAll)
			interfaceMap[iface] = true
		}
	}
	for _, _interface := range interfaces {
		wanted, ok := interfaceMap[_interface.Name]
		if wanted && ok { // Simple case
			totalBytesRecv += _interface.BytesRecv
			totalBytesSent += _interface.BytesSent
		} else if ok { // Present but unwanted
			continue
		} else if interfaceMap[NetInterfaceAll] { // Capture other
			totalBytesRecv += _interface.BytesRecv
			totalBytesSent += _interface.BytesSent
		}
	}

	var recentBytesRecv uint64
	var recentBytesSent uint64

	if net.totalBytesRecv != 0 { // if this isn't the first update
		recentBytesRecv = totalBytesRecv - net.totalBytesRecv
		recentBytesSent = totalBytesSent - net.totalBytesSent

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

		net.Lines[0].Data = append(net.Lines[0].Data, int(recentBytesRecv))
		net.Lines[1].Data = append(net.Lines[1].Data, int(recentBytesSent))
		if net.sentMetric != nil {
			net.sentMetric.Add(float64(recentBytesSent))
			net.recvMetric.Add(float64(recentBytesRecv))
		}
	}

	// used in later calls to update
	net.totalBytesRecv = totalBytesRecv
	net.totalBytesSent = totalBytesSent

	rx, tx := "RX/s", "TX/s"
	if net.Mbps {
		rx, tx = "mbps", "mbps"
	}
	format := " %s: %9.1f %2s/s"

	var total, recent uint64
	var label, unitRecent, rate string
	var recentConverted float64
	// render widget titles
	for i := 0; i < 2; i++ {
		if i == 0 {
			total, label, rate, recent = totalBytesRecv, "RX", rx, recentBytesRecv
		} else {
			total, label, rate, recent = totalBytesSent, "TX", tx, recentBytesSent
		}

		totalConverted, unitTotal := utils.ConvertBytes(total)
		if net.Mbps {
			recentConverted, unitRecent, format = float64(recent)*0.000008, "", " %s: %11.3f %2s"
		} else {
			recentConverted, unitRecent = utils.ConvertBytes(recent)
		}

		net.Lines[i].Title1 = fmt.Sprintf(" Total %s: %5.1f %s", label, totalConverted, unitTotal)
		net.Lines[i].Title2 = fmt.Sprintf(format, rate, recentConverted, unitRecent)
	}
}
