package widgets

import (
	"image"
	"log"
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

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("could not get hostname: %v", err)
		return
	}
	buf.SetString(
		hostname,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(self.Inner.Min.X, self.Inner.Min.Y+(self.Inner.Dy()/2)),
	)

	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04:05")
	buf.SetString(
		formattedTime,
		ui.NewStyle(ui.ColorWhite),
		image.Pt(
			self.Inner.Min.X+(self.Inner.Dx()/2)-len(formattedTime)/2,
			self.Inner.Min.Y+(self.Inner.Dy()/2),
		),
	)

	buf.SetString(
		"gotop",
		ui.NewStyle(ui.ColorWhite),
		image.Pt(
			self.Inner.Max.X-6,
			self.Inner.Min.Y+(self.Inner.Dy()/2),
		),
	)
}
