package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"time"

	psCPU "github.com/shirou/gopsutil/cpu"

	ui "github.com/cjbassi/gotop/src/termui"
	"github.com/cjbassi/gotop/src/utils"
)

const (
	UP_ARROW   = "▲"
	DOWN_ARROW = "▼"
)

type ProcSortMethod string

const (
	ProcSortCpu ProcSortMethod = "c"
	ProcSortMem                = "m"
	ProcSortPid                = "p"
)

type Proc struct {
	Pid         int
	CommandName string
	FullCommand string
	Cpu         float64
	Mem         float64
}

type ProcWidget struct {
	*ui.Table
	cpuCount         int
	updateInterval   time.Duration
	sortMethod       ProcSortMethod
	groupedProcs     []Proc
	ungroupedProcs   []Proc
	showGroupedProcs bool
}

func NewProcWidget() *ProcWidget {
	cpuCount, err := psCPU.Counts(false)
	if err != nil {
		log.Printf("failed to get CPU count from gopsutil: %v", err)
	}
	self := &ProcWidget{
		Table:            ui.NewTable(),
		updateInterval:   time.Second,
		cpuCount:         cpuCount,
		sortMethod:       ProcSortCpu,
		showGroupedProcs: true,
	}
	self.Title = " Processes "
	self.ShowCursor = true
	self.ShowLocation = true
	self.ColGap = 3
	self.PadLeft = 2
	self.ColResizer = func() {
		self.ColWidths = []int{
			5, utils.MaxInt(self.Inner.Dx()-26, 10), 4, 4,
		}
	}

	self.UniqueCol = 0
	if self.showGroupedProcs {
		self.UniqueCol = 1
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

func (self *ProcWidget) update() {
	procs, err := getProcs()
	if err != nil {
		log.Printf("failed to retrieve processes: %v", err)
		return
	}

	// have to iterate over the entry number in order to modify the array in place
	for i := range procs {
		procs[i].Cpu /= float64(self.cpuCount)
	}

	self.ungroupedProcs = procs
	self.groupedProcs = groupProcs(procs)

	self.sortProcs()
	self.convertProcsToTableRows()
}

// sortProcs sorts either the grouped or ungrouped []Process based on the sortMethod.
// Called with every update, when the sort method is changed, and when processes are grouped and ungrouped.
func (self *ProcWidget) sortProcs() {
	self.Header = []string{"Count", "Command", "CPU%", "Mem%"}

	if !self.showGroupedProcs {
		self.Header[0] = "PID"
	}

	var procs *[]Proc
	if self.showGroupedProcs {
		procs = &self.groupedProcs
	} else {
		procs = &self.ungroupedProcs
	}

	switch self.sortMethod {
	case ProcSortCpu:
		sort.Sort(sort.Reverse(SortProcsByCpu(*procs)))
		self.Header[2] += DOWN_ARROW
	case ProcSortPid:
		if self.showGroupedProcs {
			sort.Sort(sort.Reverse(SortProcsByPid(*procs)))
		} else {
			sort.Sort(SortProcsByPid(*procs))
		}
		self.Header[0] += DOWN_ARROW
	case ProcSortMem:
		sort.Sort(sort.Reverse(SortProcsByMem(*procs)))
		self.Header[3] += DOWN_ARROW
	}
}

// convertProcsToTableRows converts a []Proc to a [][]string and sets it to the table Rows
func (self *ProcWidget) convertProcsToTableRows() {
	var procs *[]Proc
	if self.showGroupedProcs {
		procs = &self.groupedProcs
	} else {
		procs = &self.ungroupedProcs
	}
	strings := make([][]string, len(*procs))
	for i := range *procs {
		strings[i] = make([]string, 4)
		strings[i][0] = strconv.Itoa(int((*procs)[i].Pid))
		if self.showGroupedProcs {
			strings[i][1] = (*procs)[i].CommandName
		} else {
			strings[i][1] = (*procs)[i].FullCommand
		}
		strings[i][2] = fmt.Sprintf("%4s", strconv.FormatFloat((*procs)[i].Cpu, 'f', 1, 64))
		strings[i][3] = fmt.Sprintf("%4s", strconv.FormatFloat(float64((*procs)[i].Mem), 'f', 1, 64))
	}
	self.Rows = strings
}

func (self *ProcWidget) ChangeProcSortMethod(method ProcSortMethod) {
	if self.sortMethod != method {
		self.sortMethod = method
		self.ScrollTop()
		self.sortProcs()
		self.convertProcsToTableRows()
	}
}

func (self *ProcWidget) ToggleShowingGroupedProcs() {
	self.showGroupedProcs = !self.showGroupedProcs
	if self.showGroupedProcs {
		self.UniqueCol = 1
	} else {
		self.UniqueCol = 0
	}
	self.ScrollTop()
	self.sortProcs()
	self.convertProcsToTableRows()
}

// KillProc kills a process or group of processes depending on if we're displaying the processes grouped or not.
func (self *ProcWidget) KillProc() {
	self.SelectedItem = ""
	command := "kill"
	if self.UniqueCol == 1 {
		command = "pkill"
	}
	cmd := exec.Command(command, self.Rows[self.SelectedRow][self.UniqueCol])
	cmd.Start()
	cmd.Wait()
}

// groupProcs groupes a []Proc based on command name.
// The first field changes from PID to count.
// Cpu and Mem are added together for each Proc.
func groupProcs(procs []Proc) []Proc {
	groupedProcsMap := make(map[string]Proc)
	for _, proc := range procs {
		val, ok := groupedProcsMap[proc.CommandName]
		if ok {
			groupedProcsMap[proc.CommandName] = Proc{
				val.Pid + 1,
				val.CommandName,
				"",
				val.Cpu + proc.Cpu,
				val.Mem + proc.Mem,
			}
		} else {
			groupedProcsMap[proc.CommandName] = Proc{
				1,
				proc.CommandName,
				"",
				proc.Cpu,
				proc.Mem,
			}
		}
	}

	groupedProcsList := make([]Proc, len(groupedProcsMap))
	i := 0
	for _, val := range groupedProcsMap {
		groupedProcsList[i] = val
		i++
	}

	return groupedProcsList
}

// []Proc Sorting //////////////////////////////////////////////////////////////

type SortProcsByCpu []Proc

// Len implements Sort interface
func (self SortProcsByCpu) Len() int {
	return len(self)
}

// Swap implements Sort interface
func (self SortProcsByCpu) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// Less implements Sort interface
func (self SortProcsByCpu) Less(i, j int) bool {
	return self[i].Cpu < self[j].Cpu
}

type SortProcsByPid []Proc

// Len implements Sort interface
func (self SortProcsByPid) Len() int {
	return len(self)
}

// Swap implements Sort interface
func (self SortProcsByPid) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// Less implements Sort interface
func (self SortProcsByPid) Less(i, j int) bool {
	return self[i].Pid < self[j].Pid
}

type SortProcsByMem []Proc

// Len implements Sort interface
func (self SortProcsByMem) Len() int {
	return len(self)
}

// Swap implements Sort interface
func (self SortProcsByMem) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// Less implements Sort interface
func (self SortProcsByMem) Less(i, j int) bool {
	return self[i].Mem < self[j].Mem
}
