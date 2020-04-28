package widgets

import (
	"image"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
)

type StatusBar struct {
	ui.Block
}

func NewStatusBar() *StatusBar {
	self := &StatusBar{*ui.NewBlock()}
	self.Border = false
	return self
}

func (sb *StatusBar) Draw(buf *ui.Buffer) {
	sb.Block.Draw(buf)

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("could not get hostname: %v", err)
		return
	}
	buf.SetString(
		hostname,
		ui.Theme.Default,
		image.Pt(sb.Inner.Min.X, sb.Inner.Min.Y+(sb.Inner.Dy()/2)),
	)

	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04:05")
	buf.SetString(
		formattedTime,
		ui.Theme.Default,
		image.Pt(
			sb.Inner.Min.X+(sb.Inner.Dx()/2)-len(formattedTime)/2,
			sb.Inner.Min.Y+(sb.Inner.Dy()/2),
		),
	)

	buf.SetString(
		"gotop",
		ui.Theme.Default,
		image.Pt(
			sb.Inner.Max.X-6,
			sb.Inner.Min.Y+(sb.Inner.Dy()/2),
		),
	)
}
