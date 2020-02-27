package termui

import (
	. "github.com/gizak/termui/v3"
	gizak "github.com/gizak/termui/v3/widgets"
)

// LineGraph implements a line graph of data points.
type Gauge struct {
	*gizak.Gauge
}

func NewGauge() *Gauge {
	return &Gauge{
		Gauge: gizak.NewGauge(),
	}
}

func (self *Gauge) Draw(buf *Buffer) {
	self.Gauge.Draw(buf)
	self.Gauge.SetRect(self.Min.X, self.Min.Y, self.Inner.Dx(), self.Inner.Dy())
}
