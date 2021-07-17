package devices

import "log"

var memFuncs []func(map[string]MemoryInfo) map[string]error

// TODO Colors are wrong for #mem > 2
type MemoryInfo struct {
	Total       uint64
	Used        uint64
	UsedPercent float64
}

func RegisterMem(f func(map[string]MemoryInfo) map[string]error) {
	memFuncs = append(memFuncs, f)
}

func UpdateMem(mem map[string]MemoryInfo) {
	for _, f := range memFuncs {
		errs := f(mem)
		if errs != nil {
			for k, e := range errs {
				log.Printf("%s: %s", k, e)
			}
		}
	}
}
