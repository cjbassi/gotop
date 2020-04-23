package widgets

import (
	"fmt"
	"image"
	"sort"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xxxserxxx/gotop/v4/devices"
	"github.com/xxxserxxx/gotop/v4/utils"
)

type TempScale rune

const (
	Celsius    TempScale = 'C'
	Fahrenheit           = 'F'
)

type TempWidget struct {
	*ui.Block      // inherits from Block instead of a premade Widget
	updateInterval time.Duration
	Data           map[string]int
	TempThreshold  int
	TempLowColor   ui.Color
	TempHighColor  ui.Color
	TempScale      TempScale
	tempsMetric    map[string]prometheus.Gauge
}

// TODO: state:deferred 156 Added temperatures for NVidia GPUs (azak-azkaran/master). Crashes on non-nvidia machines.
func NewTempWidget(tempScale TempScale, filter []string) *TempWidget {
	self := &TempWidget{
		Block:          ui.NewBlock(),
		updateInterval: time.Second * 5,
		Data:           make(map[string]int),
		TempThreshold:  80,
		TempScale:      tempScale,
	}
	self.Title = " Temperatures "
	if len(filter) > 0 {
		for _, t := range filter {
			self.Data[t] = 0
		}
	} else {
		for _, t := range devices.Devices(devices.Temperatures, false) {
			self.Data[t] = 0
		}
	}

	if tempScale == Fahrenheit {
		self.TempThreshold = utils.CelsiusToFahrenheit(self.TempThreshold)
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

func (self *TempWidget) EnableMetric() {
	self.tempsMetric = make(map[string]prometheus.Gauge)
	for k, v := range self.Data {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gotop",
			Subsystem: "temp",
			Name:      k,
		})
		gauge.Set(float64(v))
		prometheus.MustRegister(gauge)
		self.tempsMetric[k] = gauge
	}
}

// Custom Draw method instead of inheriting from a generic Widget.
func (self *TempWidget) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	var keys []string
	for key := range self.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for y, key := range keys {
		if y+1 > self.Inner.Dy() {
			break
		}

		var fg ui.Color
		if self.Data[key] < self.TempThreshold {
			fg = self.TempLowColor
		} else {
			fg = self.TempHighColor
		}

		s := ui.TrimString(key, (self.Inner.Dx() - 4))
		buf.SetString(s,
			ui.Theme.Default,
			image.Pt(self.Inner.Min.X, self.Inner.Min.Y+y),
		)

		if self.tempsMetric != nil {
			self.tempsMetric[key].Set(float64(self.Data[key]))
		}
		temperature := fmt.Sprintf("%3d°%c", self.Data[key], self.TempScale)

		buf.SetString(
			temperature,
			ui.NewStyle(fg),
			image.Pt(self.Inner.Max.X-(len(temperature)-1), self.Inner.Min.Y+y),
		)
	}
}

func (self *TempWidget) update() {
	devices.UpdateTemps(self.Data)
	for name, val := range self.Data {
		if self.TempScale == Fahrenheit {
			self.Data[name] = utils.CelsiusToFahrenheit(val)
		} else {
			self.Data[name] = val
		}
	}
}
