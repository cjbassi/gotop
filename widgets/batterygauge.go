package widgets

import (
	"fmt"
	"log"

	"time"

	"github.com/distatus/battery"
	"github.com/prometheus/client_golang/prometheus"

	. "github.com/xxxserxxx/gotop/v3/termui"
)

// FIXME 3.5.1 is 0% always
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
		log.Printf("error setting up batteries: %v", err)
		return
	}
	mx := 0.0
	cu := 0.0
	charging := "%d%% ⚡%s"
	rate := 0.0
	for _, bat := range bats {
		mx += bat.Full
		cu += bat.Current
		if rate < bat.ChargeRate {
			rate = bat.ChargeRate
		}
		if bat.State == battery.Charging {
			charging = "%d%% 🔌%s"
		}
	}
	tn := (mx - cu) / rate
	d, _ := time.ParseDuration(fmt.Sprintf("%fh", tn))
	self.Percent = int((cu / mx) * 100.0)
	self.Label = fmt.Sprintf(charging, self.Percent, d.Truncate(time.Minute))
	if self.metric != nil {
		self.metric.Set(cu / mx)
	}
}
