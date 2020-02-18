package widgets

import (
	"fmt"
	"image"
	"sort"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/xxxserxxx/gotop/utils"
)

type TempScale int

const (
	Celsius    TempScale = 0
	Fahrenheit           = 1
	Disabled             = 2
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

func NewTempWidget(tempScale TempScale) *TempWidget {
	self := &TempWidget{
		Block:          ui.NewBlock(),
		updateInterval: time.Second * 5,
		Data:           make(map[string]int),
		TempThreshold:  80,
		TempScale:      tempScale,
	}
	self.Title = " Temperatures "

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

		switch self.TempScale {
		case Fahrenheit:
			buf.SetString(
				fmt.Sprintf("%3dF", self.Data[key]),
				ui.NewStyle(fg),
				image.Pt(self.Inner.Max.X-4, self.Inner.Min.Y+y),
			)
		case Celsius:
			buf.SetString(
				fmt.Sprintf("%3dC", self.Data[key]),
				ui.NewStyle(fg),
				image.Pt(self.Inner.Max.X-4, self.Inner.Min.Y+y),
			)
		}
		self.tempsMetric[key].Set(float64(self.Data[key]))
	}
}
