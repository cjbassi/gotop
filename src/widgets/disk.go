package widgets

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psDisk "github.com/shirou/gopsutil/disk"
)

type Partition struct {
	Device      string
	Mount       string
	TotalRead   uint64
	TotalWrite  uint64
	CurRead     string
	CurWrite    string
	UsedPercent int
	Free        string
}

type Disk struct {
	*ui.Table
	interval   time.Duration
	Partitions map[string]*Partition
}

func NewDisk() *Disk {
	self := &Disk{
		Table:      ui.NewTable(),
		interval:   time.Second,
		Partitions: make(map[string]*Partition),
	}
	self.Label = "Disk Usage"
	self.Header = []string{"Disk", "Mount", "Used", "Free", "R/s", "W/s"}
	self.Gap = 2
	self.ColResizer = self.ColResize

	self.update()

	go func() {
		ticker := time.NewTicker(self.interval)
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Disk) update() {
	Partitions, err := psDisk.Partitions(false)
	if err != nil {
		log.Printf("failed to get disk partitions from gopsutil: %v", err)
	}

	// add partition if it's new
	for _, Part := range Partitions {
		// don't show loop devices
		if strings.HasPrefix(Part.Device, "/dev/loop") {
			continue
		}
		// don't show docker container filesystems
		if strings.HasPrefix(Part.Mountpoint, "/var/lib/docker/") {
			continue
		}
		// check if partition doesn't already exist in our list
		if _, ok := self.Partitions[Part.Device]; !ok {
			self.Partitions[Part.Device] = &Partition{
				Device: Part.Device,
				Mount:  Part.Mountpoint,
			}
		}
	}

	// delete a partition if it no longer exists
	todelete := []string{}
	for key := range self.Partitions {
		exists := false
		for _, Part := range Partitions {
			if key == Part.Device {
				exists = true
				break
			}
		}
		if !exists {
			todelete = append(todelete, key)
		}
	}
	for _, val := range todelete {
		delete(self.Partitions, val)
	}

	// updates partition info
	for _, Part := range self.Partitions {
		usage, err := psDisk.Usage(Part.Mount)
		if err != nil {
			log.Printf("failed to get partition usage statistics from gopsutil: %v. Part: %v", err, Part)
			continue
		}
		Part.UsedPercent = int(usage.UsedPercent)

		Free, Mag := utils.ConvertBytes(usage.Free)
		Part.Free = fmt.Sprintf("%3d%s", uint64(Free), Mag)

		ret, err := psDisk.IOCounters(Part.Device)
		if err != nil {
			log.Printf("failed to get partition read/write info from gopsutil: %v. Part: %v", err, Part)
			continue
		}
		data := ret[Part.Device]
		curRead, curWrite := data.ReadBytes, data.WriteBytes
		if Part.TotalRead != 0 { // if this isn't the first update
			readRecent := curRead - Part.TotalRead
			writeRecent := curWrite - Part.TotalWrite

			readFloat, unitRead := utils.ConvertBytes(readRecent)
			writeFloat, unitWrite := utils.ConvertBytes(writeRecent)
			readRecent, writeRecent = uint64(readFloat), uint64(writeFloat)
			Part.CurRead = fmt.Sprintf("%d%s", readRecent, unitRead)
			Part.CurWrite = fmt.Sprintf("%d%s", writeRecent, unitWrite)
		} else {
			Part.CurRead = fmt.Sprintf("%d%s", 0, "B")
			Part.CurWrite = fmt.Sprintf("%d%s", 0, "B")
		}
		Part.TotalRead, Part.TotalWrite = curRead, curWrite
	}

	// converts self.Partitions into self.Rows which is a [][]String
	sortedPartitions := []string{}
	for seriesName := range self.Partitions {
		sortedPartitions = append(sortedPartitions, seriesName)
	}
	sort.Strings(sortedPartitions)

	self.Rows = make([][]string, len(self.Partitions))
	for i, key := range sortedPartitions {
		Part := self.Partitions[key]
		self.Rows[i] = make([]string, 6)
		self.Rows[i][0] = strings.Replace(strings.Replace(Part.Device, "/dev/", "", -1), "mapper/", "", -1)
		self.Rows[i][1] = Part.Mount
		self.Rows[i][2] = fmt.Sprintf("%d%%", Part.UsedPercent)
		self.Rows[i][3] = Part.Free
		self.Rows[i][4] = Part.CurRead
		self.Rows[i][5] = Part.CurWrite
	}
}

// ColResize overrides the default ColResize in the termui table.
func (self *Disk) ColResize() {
	self.ColWidths = []int{
		utils.Max(4, (self.X-29)/2),
		utils.Max(5, (self.X-29)/2),
		4, 5, 5, 5,
	}

	self.CellXPos = []int{}
	cur := 1
	for _, w := range self.ColWidths {
		self.CellXPos = append(self.CellXPos, cur)
		cur += w
		cur += self.Gap
	}
}
