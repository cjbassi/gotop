// +build darwin

package devices

import smc "github.com/xxxserxxx/iSMC"

func init() {
	RegisterTemp(update)
	RegisterDeviceList(Temperatures, devs)
	ts = make(map[string]float32)
}

var ts map[string]float32

func update(temps map[string]int) map[string]error {
	err := smc.GetTemp(ts)
	if err != nil {
		return map[string]error{"temps": err}
	}
	for k, v := range ts {
		if _, ok := temps[k]; ok {
			temps[k] = int(v + 0.5)
		}
	}
	return nil
}

// TODO: Set reasonable default devices
// CPU (TC[01]P), GPU (TG0P), Memory (Ts0S) and Disk (TH0P)
func devs() []string {
	rv := make([]string, len(smc.AppleTemp))
	for i, v := range smc.AppleTemp {
		rv[i] = v.Desc
	}
	return rv
}
