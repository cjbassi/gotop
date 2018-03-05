package widgets

import (
	"os/exec"
	"sort"
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/termui"
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
	cpuCount       int
	interval       time.Duration
	sortMethod     string
	groupedProcs   []Process
	ungroupedProcs []Process
	group          bool
	KeyPressed     chan bool
}

func NewProc(loaded, keyPressed chan bool) *Proc {
	cpuCount, _ := psCPU.Counts(false)
	p := &Proc{
		Table:      ui.NewTable(),
		interval:   time.Second,
		cpuCount:   cpuCount,
		sortMethod: "c",
		group:      true,
		KeyPressed: keyPressed,
	}
	p.ColResizer = p.ColResize
	p.Label = "Process List"
	p.ColWidths = []int{5, 10, 4, 4}

	p.UniqueCol = 0
	if p.group {
		p.UniqueCol = 1
	}

	p.keyBinds()

	go func() {
		p.update()
		loaded <- true
	}()

	ticker := time.NewTicker(p.interval)
	go func() {
		for range ticker.C {
			p.update()
		}
	}()

	return p
}

func (p *Proc) update() {
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
			cpu / float64(p.cpuCount),
			mem,
		}
	}
	p.ungroupedProcs = processes
	p.groupedProcs = Group(processes)

	p.Sort()
}

// Sort sorts either the grouped or ungrouped []Process based on the sortMethod.
// Called with every update, when the sort method is changed, and when processes are grouped and ungrouped.
func (p *Proc) Sort() {
	p.Header = []string{"Count", "Command", "CPU%", "Mem%"}

	if !p.group {
		p.Header[0] = "PID"
	}

	processes := &p.ungroupedProcs
	if p.group {
		processes = &p.groupedProcs
	}

	switch p.sortMethod {
	case "c":
		sort.Sort(sort.Reverse(ProcessByCPU(*processes)))
		p.Header[2] += DOWN
	case "p":
		if p.group {
			sort.Sort(sort.Reverse(ProcessByPID(*processes)))
		} else {
			sort.Sort(ProcessByPID(*processes))
		}
		p.Header[0] += DOWN
	case "m":
		sort.Sort(sort.Reverse(ProcessByMem(*processes)))
		p.Header[3] += DOWN
	}

	p.Rows = FieldsToStrings(*processes)
}

// ColResize overrides the default ColResize in the termui table.
func (p *Proc) ColResize() {
	// calculate gap size based on total width
	p.Gap = 3
	if p.X < 50 {
		p.Gap = 1
	} else if p.X < 75 {
		p.Gap = 2
	}

	p.CellXPos = []int{
		p.Gap,
		p.Gap + p.ColWidths[0] + p.Gap,
		p.X - p.Gap - p.ColWidths[3] - p.Gap - p.ColWidths[2],
		p.X - p.Gap - p.ColWidths[3],
	}

	rowWidth := p.Gap + p.ColWidths[0] + p.Gap + p.ColWidths[1] + p.Gap + p.ColWidths[2] + p.Gap + p.ColWidths[3] + p.Gap

	// only renders a column if it fits
	if p.X < (rowWidth - p.Gap - p.ColWidths[3]) {
		p.ColWidths[2] = 0
		p.ColWidths[3] = 0
	} else if p.X < rowWidth {
		p.CellXPos[2] = p.CellXPos[3]
		p.ColWidths[3] = 0
	}
}

func (p *Proc) keyBinds() {
	ui.On("MouseLeft", func(e ui.Event) {
		p.Click(e.MouseX, e.MouseY)
		p.KeyPressed <- true
	})

	ui.On("MouseWheelUp", "MouseWheelDown", func(e ui.Event) {
		switch e.Key {
		case "MouseWheelDown":
			p.Down()
		case "MouseWheelUp":
			p.Up()
		}
		p.KeyPressed <- true
	})

	ui.On("<up>", "<down>", func(e ui.Event) {
		switch e.Key {
		case "<up>":
			p.Up()
		case "<down>":
			p.Down()
		}
		p.KeyPressed <- true
	})

	viKeys := []string{"j", "k", "gg", "G", "C-d", "C-u", "C-f", "C-b"}
	ui.On(viKeys, func(e ui.Event) {
		switch e.Key {
		case "j":
			p.Down()
		case "k":
			p.Up()
		case "gg":
			p.Top()
		case "G":
			p.Bottom()
		case "C-d":
			p.HalfPageDown()
		case "C-u":
			p.HalfPageUp()
		case "C-f":
			p.PageDown()
		case "C-b":
			p.PageUp()
		}
		p.KeyPressed <- true
	})

	ui.On("dd", func(e ui.Event) {
		p.Kill()
	})

	ui.On("<tab>", func(e ui.Event) {
		p.group = !p.group
		if p.group {
			p.UniqueCol = 1
		} else {
			p.UniqueCol = 0
		}
		p.sortMethod = "c"
		p.Sort()
		p.Top()
		p.KeyPressed <- true
	})

	ui.On("m", "c", "p", func(e ui.Event) {
		if p.sortMethod != e.Key {
			p.sortMethod = e.Key
			p.Top()
			p.Sort()
			p.KeyPressed <- true
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
		strings[i][2] = strconv.FormatFloat(p.CPU, 'f', 1, 64)
		strings[i][3] = strconv.FormatFloat(float64(p.Mem), 'f', 1, 32)
	}
	return strings
}

// Kill kills process or group of processes.
func (p *Proc) Kill() {
	p.SelectedItem = ""
	command := "kill"
	if p.UniqueCol == 1 {
		command = "pkill"
	}
	cmd := exec.Command(command, p.Rows[p.SelectedRow][p.UniqueCol])
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
