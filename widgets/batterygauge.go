package widgets

import (
	"fmt"
	"log"
	//"math"
	//"strconv"
	"time"

	"github.com/distatus/battery"
	"github.com/prometheus/client_golang/prometheus"

	. "github.com/xxxserxxx/gotop/termui"
)

type BatteryGauge struct {
	*Gauge
	metric prometheus.Gauge
}

func NewBatteryGauge() *BatteryGauge {
	self := &BatteryGauge{Gauge: NewGauge()}
	self.Title = " Power Level "

	self.update()

	go func() {
		for range time.NewTicker(time.Second).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (b *BatteryGauge) EnableMetric() {
	bats, err := battery.GetAll()
	if err != nil {
		log.Printf("error setting up metrics: %v", err)
		return
	}
	mx := 0.0
	cu := 0.0
	for _, bat := range bats {
		mx += bat.Full
		cu += bat.Current
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gotop",
			Subsystem: "battery",
			Name:      "total",
		})
		gauge.Set(cu / mx)
		b.metric = gauge
		prometheus.MustRegister(gauge)
	}
}

func (self *BatteryGauge) update() {
	bats, err := battery.GetAll()
	if err != nil {
		log.Printf("error setting up metrics: %v", err)
		return
	}
	mx := 0.0
	cu := 0.0
	for _, bat := range bats {
		mx += bat.Full
		cu += bat.Current
	}
	self.Percent = int((cu / mx) * 100.0)
	self.Label = fmt.Sprintf("%d%%", self.Percent)
	if self.metric != nil {
		self.metric.Set(cu / mx)
	}
}
