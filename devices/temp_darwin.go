// +build darwin

package devices

import smc "github.com/xxxserxxx/iSMC"

func init() {
	RegisterTemp(update)
	ts = make(map[string]float32)
}

var ts map[string]float32

func update(temps map[string]int) map[string]error {
	err := smc.GetTemp(ts)
	if err != nil {
		return map[string]error{"temps": err}
	}
	for k, v := range ts {
		temps[k] = int(v + 0.5)
	}
	return nil
}
