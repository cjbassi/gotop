package widgets

import (
	"fmt"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psDisk "github.com/shirou/gopsutil/disk"
)

type Disk struct {
	*ui.Gauge
	fs       string // which filesystem to get the disk usage of
	interval time.Duration
}

func NewDisk() *Disk {
	self := &Disk{
		Gauge:    ui.NewGauge(),
		fs:       "/",
		interval: time.Second * 5,
	}
	self.Label = "Disk Usage"

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Disk) update() {
	usage, _ := psDisk.Usage(self.fs)
	self.Percent = int(usage.UsedPercent)
	self.Description = fmt.Sprintf(" (%dGB free)", int(utils.BytesToGB(usage.Free)))
}
