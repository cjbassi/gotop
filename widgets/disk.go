package widgets

import (
	"fmt"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	"github.com/cjbassi/gotop/utils"
	disk "github.com/shirou/gopsutil/disk"
)

type Disk struct {
	*ui.Gauge
	fs       string // which filesystem to get the disk usage of
	interval time.Duration
}

func NewDisk() *Disk {
	// get root filesystem usage
	d := &Disk{ui.NewGauge(), "/", time.Second * 5}
	d.Label = "Disk Usage"

	go d.update()
	ticker := time.NewTicker(d.interval)
	go func() {
		for range ticker.C {
			d.update()
		}
	}()

	return d
}

func (d *Disk) update() {
	disk, _ := disk.Usage(d.fs)
	d.Percent = int(disk.UsedPercent)
	d.Description = fmt.Sprintf(" (%dGB free)", int(utils.BytesToGB(disk.Free)))
}
