package widgets

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	tui "github.com/gizak/termui/v3"
	"github.com/xxxserxxx/gotop/v4/devices"
	ui "github.com/xxxserxxx/gotop/v4/termui"
	"github.com/xxxserxxx/gotop/v4/utils"
)

const (
	_downArrow = "▼"
)

type ProcSortMethod string

const (
	ProcSortCPU ProcSortMethod = "c"
	ProcSortMem                = "m"
	ProcSortPid                = "p"
)

type Proc struct {
	Pid         int
	CommandName string
	FullCommand string
	CPU         float64
	Mem         float64
}

type ProcWidget struct {
	*ui.Table
	entry            *ui.Entry
	cpuCount         int
	updateInterval   time.Duration
	sortMethod       ProcSortMethod
	filter           string
	groupedProcs     []Proc
	ungroupedProcs   []Proc
	showGroupedProcs bool
}

func NewProcWidget() *ProcWidget {
	cpuCount, err := devices.CpuCount()
	if err != nil {
		log.Printf("failed to get CPU count from gopsutil: %v", err)
	}
	self := &ProcWidget{
		Table:            ui.NewTable(),
		updateInterval:   time.Second,
		cpuCount:         cpuCount,
		sortMethod:       ProcSortCPU,
		showGroupedProcs: true,
		filter:           "",
	}
	self.entry = &ui.Entry{
		Style: self.TitleStyle,
		Label: " Filter: ",
		Value: "",
		UpdateCallback: func(val string) {
			self.filter = val
			self.update()
		},
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

func (proc *ProcWidget) EnableMetric() {
	// There's (currently) no metric for this
}

func (proc *ProcWidget) SetEditingFilter(editing bool) {
	proc.entry.SetEditing(editing)
}

func (proc *ProcWidget) HandleEvent(e tui.Event) bool {
	return proc.entry.HandleEvent(e)
}

func (proc *ProcWidget) SetRect(x1, y1, x2, y2 int) {
	proc.Table.SetRect(x1, y1, x2, y2)
	proc.entry.SetRect(x1+2, y2-1, x2-2, y2)
}

func (proc *ProcWidget) Draw(buf *tui.Buffer) {
	proc.Table.Draw(buf)
	proc.entry.Draw(buf)
}

func (proc *ProcWidget) filterProcs(procs []Proc) []Proc {
	if proc.filter == "" {
		return procs
	}
	var filtered []Proc
	for _, p := range procs {
		if strings.Contains(p.FullCommand, proc.filter) || strings.Contains(fmt.Sprintf("%d", p.Pid), proc.filter) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (proc *ProcWidget) update() {
	procs, err := getProcs()
	if err != nil {
		log.Printf("failed to retrieve processes: %v", err)
		return
	}

	// have to iterate over the entry number in order to modify the array in place
	for i := range procs {
		procs[i].CPU /= float64(proc.cpuCount)
	}

	procs = proc.filterProcs(procs)
	proc.ungroupedProcs = procs
	proc.groupedProcs = groupProcs(procs)

	proc.sortProcs()
	proc.convertProcsToTableRows()
}

// sortProcs sorts either the grouped or ungrouped []Process based on the sortMethod.
// Called with every update, when the sort method is changed, and when processes are grouped and ungrouped.
func (proc *ProcWidget) sortProcs() {
	proc.Header = []string{"Count", "Command", "CPU%", "Mem%"}

	if !proc.showGroupedProcs {
		proc.Header[0] = "PID"
	}

	var procs *[]Proc
	if proc.showGroupedProcs {
		procs = &proc.groupedProcs
	} else {
		procs = &proc.ungroupedProcs
	}

	switch proc.sortMethod {
	case ProcSortCPU:
		sort.Sort(sort.Reverse(SortProcsByCPU(*procs)))
		proc.Header[2] += _downArrow
	case ProcSortPid:
		if proc.showGroupedProcs {
			sort.Sort(sort.Reverse(SortProcsByPid(*procs)))
		} else {
			sort.Sort(SortProcsByPid(*procs))
		}
		proc.Header[0] += _downArrow
	case ProcSortMem:
		sort.Sort(sort.Reverse(SortProcsByMem(*procs)))
		proc.Header[3] += _downArrow
	}
}

// convertProcsToTableRows converts a []Proc to a [][]string and sets it to the table Rows
func (proc *ProcWidget) convertProcsToTableRows() {
	var procs *[]Proc
	if proc.showGroupedProcs {
		procs = &proc.groupedProcs
	} else {
		procs = &proc.ungroupedProcs
	}
	strings := make([][]string, len(*procs))
	for i := range *procs {
		strings[i] = make([]string, 4)
		strings[i][0] = strconv.Itoa(int((*procs)[i].Pid))
		if proc.showGroupedProcs {
			strings[i][1] = (*procs)[i].CommandName
		} else {
			strings[i][1] = (*procs)[i].FullCommand
		}
		strings[i][2] = fmt.Sprintf("%4s", strconv.FormatFloat((*procs)[i].CPU, 'f', 1, 64))
		strings[i][3] = fmt.Sprintf("%4s", strconv.FormatFloat(float64((*procs)[i].Mem), 'f', 1, 64))
	}
	proc.Rows = strings
}

func (proc *ProcWidget) ChangeProcSortMethod(method ProcSortMethod) {
	if proc.sortMethod != method {
		proc.sortMethod = method
		proc.ScrollTop()
		proc.sortProcs()
		proc.convertProcsToTableRows()
	}
}

func (proc *ProcWidget) ToggleShowingGroupedProcs() {
	proc.showGroupedProcs = !proc.showGroupedProcs
	if proc.showGroupedProcs {
		proc.UniqueCol = 1
	} else {
		proc.UniqueCol = 0
	}
	proc.ScrollTop()
	proc.sortProcs()
	proc.convertProcsToTableRows()
}

// KillProc kills a process or group of processes depending on if we're
// displaying the processes grouped or not.
func (proc *ProcWidget) KillProc(sigName string) {
	proc.SelectedItem = ""
	command := "kill"
	if proc.UniqueCol == 1 {
		command = "pkill"
	}
	cmd := exec.Command(command, "--signal", sigName, proc.Rows[proc.SelectedRow][proc.UniqueCol])
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
				val.CPU + proc.CPU,
				val.Mem + proc.Mem,
			}
		} else {
			groupedProcsMap[proc.CommandName] = Proc{
				1,
				proc.CommandName,
				"",
				proc.CPU,
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

type SortProcsByCPU []Proc

// Len implements Sort interface
func (procs SortProcsByCPU) Len() int {
	return len(procs)
}

// Swap implements Sort interface
func (procs SortProcsByCPU) Swap(i, j int) {
	procs[i], procs[j] = procs[j], procs[i]
}

// Less implements Sort interface
func (procs SortProcsByCPU) Less(i, j int) bool {
	return procs[i].CPU < procs[j].CPU
}

type SortProcsByPid []Proc

// Len implements Sort interface
func (procs SortProcsByPid) Len() int {
	return len(procs)
}

// Swap implements Sort interface
func (procs SortProcsByPid) Swap(i, j int) {
	procs[i], procs[j] = procs[j], procs[i]
}

// Less implements Sort interface
func (procs SortProcsByPid) Less(i, j int) bool {
	return procs[i].Pid < procs[j].Pid
}

type SortProcsByMem []Proc

// Len implements Sort interface
func (procs SortProcsByMem) Len() int {
	return len(procs)
}

// Swap implements Sort interface
func (procs SortProcsByMem) Swap(i, j int) {
	procs[i], procs[j] = procs[j], procs[i]
}

// Less implements Sort interface
func (procs SortProcsByMem) Less(i, j int) bool {
	return procs[i].Mem < procs[j].Mem
}
