// +build darwin

package devices

// TODO gopsutil team reports this is not needed; try getting rid of this dep
import smc "github.com/xxxserxxx/iSMC"

func init() {
	RegisterTemp(update)
	RegisterDeviceList(Temperatures, devs, defs)
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

func devs() []string {
	rv := make([]string, len(smc.AppleTemp))
	for i, v := range smc.AppleTemp {
		rv[i] = v.Desc
	}
	return rv
}

func defs() []string {
	//                     CPU 0         CPU 1         GPU           Memory        Disk
	ids := map[string]bool{"TC0P": true, "TC1P": true, "TG0P": true, "Ts0S": true, "TH0P": true}
	rv := make([]string, 0, len(ids))
	for _, v := range smc.AppleTemp {
		if ids[v.Key] {
			rv = append(rv, v.Desc)
		}
	}
	return rv
}
