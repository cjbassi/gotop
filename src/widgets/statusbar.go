package widgets

import (
	"image"
	"os"
	"time"

	ui "github.com/gizak/termui"
)

type StatusBar struct {
	ui.Block
}

func NewStatusBar() *StatusBar {
	self := &StatusBar{*ui.NewBlock()}
	self.Border = false
	return self
}

func (self *StatusBar) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	hostname, _ := os.Hostname()
	buf.SetString(
		hostname,
		ui.NewStyle(7),
		image.Pt(self.Inner.Min.X, self.Inner.Min.Y+(self.Inner.Dy()/2)),
	)

	t := time.Now()
	_time := t.Format("15:04:05")
	buf.SetString(
		_time,
		ui.NewStyle(7),
		image.Pt(
			self.Inner.Min.X+(self.Inner.Dx()/2)-len(_time)/2,
			self.Inner.Min.Y+(self.Inner.Dy()/2),
		),
	)

	buf.SetString(
		"gotop",
		ui.NewStyle(7),
		image.Pt(
			self.Inner.Max.X-6,
			self.Inner.Min.Y+(self.Inner.Dy()/2),
		),
	)
}
