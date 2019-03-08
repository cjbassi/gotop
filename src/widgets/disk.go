package widgets

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	psDisk "github.com/shirou/gopsutil/disk"

	ui "github.com/cjbassi/gotop/src/termui"
	"github.com/cjbassi/gotop/src/utils"
)

type Partition struct {
	Device               string
	MountPoint           string
	BytesRead            uint64
	BytesWritten         uint64
	BytesReadRecently    string
	BytesWrittenRecently string
	UsedPercent          uint32
	Free                 string
}

type DiskWidget struct {
	*ui.Table
	updateInterval time.Duration
	Partitions     map[string]*Partition
}

func NewDiskWidget() *DiskWidget {
	self := &DiskWidget{
		Table:          ui.NewTable(),
		updateInterval: time.Second,
		Partitions:     make(map[string]*Partition),
	}
	self.Title = " Disk Usage "
	self.Header = []string{"Disk", "Mount", "Used", "Free", "R/s", "W/s"}
	self.ColGap = 2
	self.ColResizer = func() {
		self.ColWidths = []int{
			utils.MaxInt(4, (self.Inner.Dx()-29)/2),
			utils.MaxInt(5, (self.Inner.Dx()-29)/2),
			4, 5, 5, 5,
		}
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

func (self *DiskWidget) update() {
	partitions, err := psDisk.Partitions(false)
	if err != nil {
		log.Printf("failed to get disk partitions from gopsutil: %v", err)
		return
	}

	// add partition if it's new
	for _, partition := range partitions {
		// don't show loop devices
		if strings.HasPrefix(partition.Device, "/dev/loop") {
			continue
		}
		// don't show docker container filesystems
		if strings.HasPrefix(partition.Mountpoint, "/var/lib/docker/") {
			continue
		}
		// check if partition doesn't already exist in our list
		if _, ok := self.Partitions[partition.Device]; !ok {
			self.Partitions[partition.Device] = &Partition{
				Device:     partition.Device,
				MountPoint: partition.Mountpoint,
			}
		}
	}

	// delete a partition if it no longer exists
	toDelete := []string{}
	for device := range self.Partitions {
		exists := false
		for _, partition := range partitions {
			if device == partition.Device {
				exists = true
				break
			}
		}
		if !exists {
			toDelete = append(toDelete, device)
		}
	}
	for _, device := range toDelete {
		delete(self.Partitions, device)
	}

	// updates partition info
	for _, partition := range self.Partitions {
		usage, err := psDisk.Usage(partition.MountPoint)
		if err != nil {
			log.Printf("failed to get partition usage statistics from gopsutil: %v. partition: %v", err, partition)
			continue
		}
		partition.UsedPercent = uint32(usage.UsedPercent)

		bytesFree, magnitudeFree := utils.ConvertBytes(usage.Free)
		partition.Free = fmt.Sprintf("%3d%s", uint64(bytesFree), magnitudeFree)

		ioCounters, err := psDisk.IOCounters(partition.Device)
		if err != nil {
			log.Printf("failed to get partition read/write info from gopsutil: %v. partition: %v", err, partition)
			continue
		}
		ioCounter := ioCounters[strings.Replace(partition.Device, "/dev/", "", -1)]
		bytesRead, bytesWritten := ioCounter.ReadBytes, ioCounter.WriteBytes
		if partition.BytesRead != 0 { // if this isn't the first update
			bytesReadRecently := bytesRead - partition.BytesRead
			bytesWrittenRecently := bytesWritten - partition.BytesWritten

			readFloat, readMagnitude := utils.ConvertBytes(bytesReadRecently)
			writeFloat, writeMagnitude := utils.ConvertBytes(bytesWrittenRecently)
			bytesReadRecently, bytesWrittenRecently = uint64(readFloat), uint64(writeFloat)
			partition.BytesReadRecently = fmt.Sprintf("%d%s", bytesReadRecently, readMagnitude)
			partition.BytesWrittenRecently = fmt.Sprintf("%d%s", bytesWrittenRecently, writeMagnitude)
		} else {
			partition.BytesReadRecently = fmt.Sprintf("%d%s", 0, "B")
			partition.BytesWrittenRecently = fmt.Sprintf("%d%s", 0, "B")
		}
		partition.BytesRead, partition.BytesWritten = bytesRead, bytesWritten
	}

	// converts self.Partitions into self.Rows which is a [][]String

	sortedPartitions := []string{}
	for seriesName := range self.Partitions {
		sortedPartitions = append(sortedPartitions, seriesName)
	}
	sort.Strings(sortedPartitions)

	self.Rows = make([][]string, len(self.Partitions))

	for i, key := range sortedPartitions {
		partition := self.Partitions[key]
		self.Rows[i] = make([]string, 6)
		self.Rows[i][0] = strings.Replace(strings.Replace(partition.Device, "/dev/", "", -1), "mapper/", "", -1)
		self.Rows[i][1] = partition.MountPoint
		self.Rows[i][2] = fmt.Sprintf("%d%%", partition.UsedPercent)
		self.Rows[i][3] = partition.Free
		self.Rows[i][4] = partition.BytesReadRecently
		self.Rows[i][5] = partition.BytesWrittenRecently
	}
}
