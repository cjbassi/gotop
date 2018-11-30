package widgets

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"time"

	"github.com/cjbassi/gotop/src/utils"
	ui "github.com/cjbassi/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
)

const (
	UP   = "▲"
	DOWN = "▼"
)

// Process represents each process.
type Process struct {
	PID     int
	Command string
	CPU     float64
	Mem     float64
	Args    string
}

type Proc struct {
	*ui.Table
	cpuCount       float64
	interval       time.Duration
	sortMethod     string
	groupedProcs   []Process
	ungroupedProcs []Process
	group          bool
}

func NewProc() *Proc {
	cpuCount, _ := psCPU.Counts(false)
	self := &Proc{
		Table:      ui.NewTable(),
		interval:   time.Second,
		cpuCount:   float64(cpuCount),
		sortMethod: "c",
		group:      true,
	}
	self.Label = "Processes"
	self.ColResizer = self.ColResize
	self.Cursor = true
	self.Gap = 3
	self.PadLeft = 2

	self.UniqueCol = 0
	if self.group {
		self.UniqueCol = 1
	}

	self.update()

	go func() {
		ticker := time.NewTicker(self.interval)
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

// Sort sorts either the grouped or ungrouped []Process based on the sortMethod.
// Called with every update, when the sort method is changed, and when processes are grouped and ungrouped.
func (self *Proc) Sort() {
	self.Header = []string{"Count", "Command", "CPU%", "Mem%"}

	if !self.group {
		self.Header[0] = "PID"
	}

	processes := &self.ungroupedProcs
	if self.group {
		processes = &self.groupedProcs
	}

	switch self.sortMethod {
	case "c":
		sort.Sort(sort.Reverse(ProcessByCPU(*processes)))
		self.Header[2] += DOWN
	case "p":
		if self.group {
			sort.Sort(sort.Reverse(ProcessByPID(*processes)))
		} else {
			sort.Sort(ProcessByPID(*processes))
		}
		self.Header[0] += DOWN
	case "m":
		sort.Sort(sort.Reverse(ProcessByMem(*processes)))
		self.Header[3] += DOWN
	}

	self.Rows = FieldsToStrings(*processes, self.group)
}

// ColResize overrides the default ColResize in the termui table.
func (self *Proc) ColResize() {
	self.ColWidths = []int{
		5, utils.Max(self.X-26, 10), 4, 4,
	}
}

func (self *Proc) ChangeSort(e ui.Event) {
	if self.sortMethod != e.ID {
		self.sortMethod = e.ID
		self.Top()
		self.Sort()
	}
}

func (self *Proc) Tab() {
	self.group = !self.group
	if self.group {
		self.UniqueCol = 1
	} else {
		self.UniqueCol = 0
	}
	self.Sort()
	self.Top()
}

// Group groupes a []Process based on command name.
// The first field changes from PID to count.
// CPU and Mem are added together for each Process.
func Group(P []Process) []Process {
	groupedP := make(map[string]Process)
	for _, process := range P {
		val, ok := groupedP[process.Command]
		if ok {
			groupedP[process.Command] = Process{
				val.PID + 1,
				val.Command,
				val.CPU + process.CPU,
				val.Mem + process.Mem,
				"",
			}
		} else {
			groupedP[process.Command] = Process{
				1,
				process.Command,
				process.CPU,
				process.Mem,
				"",
			}
		}
	}

	groupList := make([]Process, len(groupedP))
	var i int
	for _, val := range groupedP {
		groupList[i] = val
		i++
	}

	return groupList
}

// FieldsToStrings converts a []Process to a [][]string
func FieldsToStrings(P []Process, grouped bool) [][]string {
	strings := make([][]string, len(P))
	for i, p := range P {
		strings[i] = make([]string, 4)
		strings[i][0] = strconv.Itoa(int(p.PID))
		if grouped {
			strings[i][1] = p.Command
		} else {
			strings[i][1] = p.Args
		}
		strings[i][2] = fmt.Sprintf("%4s", strconv.FormatFloat(p.CPU, 'f', 1, 64))
		strings[i][3] = fmt.Sprintf("%4s", strconv.FormatFloat(float64(p.Mem), 'f', 1, 64))
	}
	return strings
}

// Kill kills process or group of processes.
func (self *Proc) Kill() {
	self.SelectedItem = ""
	command := "kill"
	if self.UniqueCol == 1 {
		command = "pkill"
	}
	cmd := exec.Command(command, self.Rows[self.SelectedRow][self.UniqueCol])
	cmd.Start()
	cmd.Wait()
}

/////////////////////////////////////////////////////////////////////////////////
//                              []Process Sorting                              //
/////////////////////////////////////////////////////////////////////////////////

type ProcessByCPU []Process

// Len implements Sort interface
func (P ProcessByCPU) Len() int {
	return len(P)
}

// Swap implements Sort interface
func (P ProcessByCPU) Swap(i, j int) {
	P[i], P[j] = P[j], P[i]
}

// Less implements Sort interface
func (P ProcessByCPU) Less(i, j int) bool {
	return P[i].CPU < P[j].CPU
}

type ProcessByPID []Process

// Len implements Sort interface
func (P ProcessByPID) Len() int {
	return len(P)
}

// Swap implements Sort interface
func (P ProcessByPID) Swap(i, j int) {
	P[i], P[j] = P[j], P[i]
}

// Less implements Sort interface
func (P ProcessByPID) Less(i, j int) bool {
	return P[i].PID < P[j].PID
}

type ProcessByMem []Process

// Len implements Sort interface
func (P ProcessByMem) Len() int {
	return len(P)
}

// Swap implements Sort interface
func (P ProcessByMem) Swap(i, j int) {
	P[i], P[j] = P[j], P[i]
}

// Less implements Sort interface
func (P ProcessByMem) Less(i, j int) bool {
	return P[i].Mem < P[j].Mem
}
