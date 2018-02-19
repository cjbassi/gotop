package widgets

import (
	"sort"
	"strconv"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	cpu "github.com/shirou/gopsutil/cpu"
	ps "github.com/shirou/gopsutil/process"
)

const (
	DOWN = "▼"
	UP   = "▲"
)

// Process represents each process
type Process struct {
	PID     int32
	Command string
	CPU     float64
	Mem     float32
}

// Proc is the widget
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

// Creates a new Proc widget
func NewProc(loaded, keyPressed chan bool) *Proc {
	cpuCount, _ := cpu.Counts(false)
	p := &Proc{
		Table:      ui.NewTable(),
		interval:   time.Second,
		cpuCount:   cpuCount,
		sortMethod: "c",
		group:      true,
		KeyPressed: keyPressed,
	}
	p.Label = "Process List"

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
	psProcs, _ := ps.Processes()
	processes := make([]Process, len(psProcs))
	for i, pr := range psProcs {
		pid := pr.Pid
		command, _ := pr.Name()
		cpu, _ := pr.CPUPercent()
		mem, _ := pr.MemoryPercent()

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

// Sort sorts either the grouped or ungrouped []Process based on the sortMethod
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

func (p *Proc) keyBinds() {
	ui.On("MouseLeft", func(e ui.Event) {
		p.Click(e.MouseX, e.MouseY)
		ui.Render(p)
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
	groupMap := make(map[string]Process)
	for _, p := range P {
		val, ok := groupMap[p.Command]
		if ok {
			newP := Process{
				val.PID + 1,
				val.Command,
				val.CPU + p.CPU,
				val.Mem + p.Mem,
			}
			groupMap[p.Command] = newP
		} else {
			newP := Process{
				1,
				p.Command,
				p.CPU,
				p.Mem,
			}
			groupMap[p.Command] = newP
		}
	}
	groupList := make([]Process, len(groupMap))
	i := 0
	for _, val := range groupMap {
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

////////////////////////////////////////////////////////////////////////////////
// Sorting

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
