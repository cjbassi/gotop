package widgets

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	psDisk "github.com/shirou/gopsutil/disk"

	ui "github.com/xxxserxxx/gotop/v4/termui"
	"github.com/xxxserxxx/gotop/v4/utils"
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
	metric         map[string]prometheus.Gauge
}

func NewDiskWidget() *DiskWidget {
	self := &DiskWidget{
		Table:          ui.NewTable(),
		updateInterval: time.Second,
		Partitions:     make(map[string]*Partition),
	}
	self.Title = tr.Value("widget.disk")
	self.Header = []string{tr.Value("disk.disk"), tr.Value("disk.mount"), tr.Value("disk.used"), tr.Value("disk.free"), tr.Value("disk.rs"), tr.Value("disk.ws")}
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

func (disk *DiskWidget) EnableMetric() {
	disk.metric = make(map[string]prometheus.Gauge)
	for key, part := range disk.Partitions {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gotop",
			Subsystem: "disk",
			Name:      strings.ReplaceAll(key, "/", ":"),
			//Name:      strings.Replace(strings.Replace(part.Device, "/dev/", "", -1), "mapper/", "", -1),
		})
		gauge.Set(float64(part.UsedPercent) / 100.0)
		prometheus.MustRegister(gauge)
		disk.metric[key] = gauge
	}
}

func (disk *DiskWidget) update() {
	partitions, err := psDisk.Partitions(false)
	if err != nil {
		log.Printf(tr.Value("error.setup", "disk-partitions", err.Error()))
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
		if _, ok := disk.Partitions[partition.Device]; !ok {
			disk.Partitions[partition.Device] = &Partition{
				Device:     partition.Device,
				MountPoint: partition.Mountpoint,
			}
		}
	}

	// delete a partition if it no longer exists
	toDelete := []string{}
	for device := range disk.Partitions {
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
		delete(disk.Partitions, device)
	}

	// updates partition info. We add 0.5 to all values to make sure the truncation rounds
	for _, partition := range disk.Partitions {
		usage, err := psDisk.Usage(partition.MountPoint)
		if err != nil {
			log.Printf(tr.Value("error.recovfetch", "partition-"+partition.MountPoint+"-usage", err.Error()))
			continue
		}
		partition.UsedPercent = uint32(usage.UsedPercent + 0.5)
		bytesFree, magnitudeFree := utils.ConvertBytes(usage.Free)
		partition.Free = fmt.Sprintf("%3d%s", uint64(bytesFree+0.5), magnitudeFree)

		ioCounters, err := psDisk.IOCounters(partition.Device)
		if err != nil {
			log.Printf(tr.Value("error.recovfetch", "partition-"+partition.Device+"-rw", err.Error()))
			continue
		}
		ioCounter := ioCounters[strings.Replace(partition.Device, "/dev/", "", -1)]
		bytesRead, bytesWritten := ioCounter.ReadBytes, ioCounter.WriteBytes
		if partition.BytesRead != 0 { // if this isn't the first update
			bytesReadRecently := bytesRead - partition.BytesRead
			bytesWrittenRecently := bytesWritten - partition.BytesWritten

			readFloat, readMagnitude := utils.ConvertBytes(bytesReadRecently)
			writeFloat, writeMagnitude := utils.ConvertBytes(bytesWrittenRecently)
			bytesReadRecently, bytesWrittenRecently = uint64(readFloat+0.5), uint64(writeFloat+0.5)
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
	for seriesName := range disk.Partitions {
		sortedPartitions = append(sortedPartitions, seriesName)
	}
	sort.Strings(sortedPartitions)

	disk.Rows = make([][]string, len(disk.Partitions))

	for i, key := range sortedPartitions {
		partition := disk.Partitions[key]
		disk.Rows[i] = make([]string, 6)
		disk.Rows[i][0] = strings.Replace(strings.Replace(partition.Device, "/dev/", "", -1), "mapper/", "", -1)
		disk.Rows[i][1] = partition.MountPoint
		disk.Rows[i][2] = fmt.Sprintf("%d%%", partition.UsedPercent)
		disk.Rows[i][3] = partition.Free
		disk.Rows[i][4] = partition.BytesReadRecently
		disk.Rows[i][5] = partition.BytesWrittenRecently
		if disk.metric != nil {
			if disk.metric[key] == nil {
				log.Printf(tr.Value("error.nometrics", "disk", key))
			} else {
				disk.metric[key].Set(float64(partition.UsedPercent) / 100.0)
			}
		}
	}
}
