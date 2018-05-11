package widgets

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"time"

	ui "github.com/cjbassi/termui"
	psCPU "github.com/shirou/gopsutil/cpu"
	psProc "github.com/shirou/gopsutil/process"
)

const (
	UP   = "▲"
	DOWN = "▼"
)

// Process represents each process.
type Process struct {
	PID     int32
	Command string
	CPU     float64
	Mem     float32
}

type Proc struct {
	*ui.Table
	cpuCount         int
	interval         time.Duration
	sortMethod       string
	groupedProcs     []Process
	ungroupedProcs   []Process
	group            bool
	KeyPressed       chan bool
	DefaultColWidths []int
}

func NewProc(keyPressed chan bool) *Proc {
	cpuCount, _ := psCPU.Counts(false)
	self := &Proc{
		Table:      ui.NewTable(),
		interval:   time.Second,
		cpuCount:   cpuCount,
		sortMethod: "c",
		group:      true,
		KeyPressed: keyPressed,
	}
	self.Label = "Processes"
	self.ColResizer = self.ColResize
	self.DefaultColWidths = []int{5, 10, 4, 4}
	self.ColWidths = make([]int, 4)

	self.UniqueCol = 0
	if self.group {
		self.UniqueCol = 1
	}

	self.keyBinds()

	self.update()

	ticker := time.NewTicker(self.interval)
	go func() {
		for range ticker.C {
			self.update()
		}
	}()

	return self
}

func (self *Proc) update() {
	psProcesses, _ := psProc.Processes()
	processes := make([]Process, len(psProcesses))
	for i, psProcess := range psProcesses {
		pid := psProcess.Pid
		command, _ := psProcess.Name()
		cpu, _ := psProcess.CPUPercent()
		mem, _ := psProcess.MemoryPercent()

		processes[i] = Process{
			pid,
			command,
			cpu / float64(self.cpuCount),
			mem,
		}
	}
	self.ungroupedProcs = processes
	self.groupedProcs = Group(processes)

	self.Sort()
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

	self.Rows = FieldsToStrings(*processes)
}

// ColResize overrides the default ColResize in the termui table.
func (self *Proc) ColResize() {
	copy(self.ColWidths, self.DefaultColWidths)

	self.Gap = 3

	self.CellXPos = []int{
		self.Gap,
		self.Gap + self.ColWidths[0] + self.Gap,
		(self.X + 2) - self.Gap - self.ColWidths[3] - self.Gap - self.ColWidths[2],
		(self.X + 2) - self.Gap - self.ColWidths[3],
	}

	rowWidth := self.Gap +
		self.ColWidths[0] + self.Gap +
		self.ColWidths[1] + self.Gap +
		self.ColWidths[2] + self.Gap +
		self.ColWidths[3] + self.Gap

	// only renders a column if it fits
	if self.X < (rowWidth - self.Gap - self.ColWidths[3]) {
		self.ColWidths[2] = 0
		self.ColWidths[3] = 0
	} else if self.X < rowWidth {
		self.CellXPos[2] = self.CellXPos[3]
		self.ColWidths[3] = 0
	}
}

func (self *Proc) keyBinds() {
	ui.On("<MouseLeft>", func(e ui.Event) {
		self.Click(e.MouseX, e.MouseY)
		self.KeyPressed <- true
	})

	ui.On("<MouseWheelUp>", "<MouseWheelDown>", func(e ui.Event) {
		switch e.Key {
		case "<MouseWheelDown>":
			self.Down()
		case "<MouseWheelUp>":
			self.Up()
		}
		self.KeyPressed <- true
	})

	ui.On("<up>", "<down>", func(e ui.Event) {
		switch e.Key {
		case "<up>":
			self.Up()
		case "<down>":
			self.Down()
		}
		self.KeyPressed <- true
	})

	viKeys := []string{"j", "k", "gg", "G", "<C-d>", "<C-u>", "<C-f>", "<C-b>"}
	ui.On(viKeys, func(e ui.Event) {
		switch e.Key {
		case "j":
			self.Down()
		case "k":
			self.Up()
		case "gg":
			self.Top()
		case "G":
			self.Bottom()
		case "<C-d>":
			self.HalfPageDown()
		case "<C-u>":
			self.HalfPageUp()
		case "<C-f>":
			self.PageDown()
		case "<C-b>":
			self.PageUp()
		}
		self.KeyPressed <- true
	})

	ui.On("dd", func(e ui.Event) {
		self.Kill()
	})

	ui.On("<tab>", func(e ui.Event) {
		self.group = !self.group
		if self.group {
			self.UniqueCol = 1
		} else {
			self.UniqueCol = 0
		}
		self.Sort()
		self.Top()
		self.KeyPressed <- true
	})

	ui.On("m", "c", "p", func(e ui.Event) {
		if self.sortMethod != e.Key {
			self.sortMethod = e.Key
			self.Top()
			self.Sort()
			self.KeyPressed <- true
		}
	})
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
			}
		} else {
			groupedP[process.Command] = Process{
				1,
				process.Command,
				process.CPU,
				process.Mem,
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
func FieldsToStrings(P []Process) [][]string {
	strings := make([][]string, len(P))
	for i, p := range P {
		strings[i] = make([]string, 4)
		strings[i][0] = strconv.Itoa(int(p.PID))
		strings[i][1] = p.Command
		strings[i][2] = fmt.Sprintf("%4s", strconv.FormatFloat(p.CPU, 'f', 1, 64))
		strings[i][3] = fmt.Sprintf("%4s", strconv.FormatFloat(float64(p.Mem), 'f', 1, 32))
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
