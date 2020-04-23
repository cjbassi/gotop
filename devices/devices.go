package devices

import "log"

const (
	Temperatures = "Temperatures" // Device domain for temperature sensors
)

var Domains []string = []string{Temperatures}
var _shutdownFuncs []func() error
var _devs map[string][]string
var _defaults map[string][]string

// RegisterShutdown stores a function to be called by gotop on exit, allowing
// extensions to properly release resources.  Extensions should register a
// shutdown function IFF the extension is using resources that need to be
// released.  The returned error will be logged, but no other action will be
// taken.
func RegisterShutdown(f func() error) {
	_shutdownFuncs = append(_shutdownFuncs, f)
}

// Shutdown will be called by the `main()` function if gotop is exited
// cleanly.  It will call all of the registered shutdown functions of devices,
// logging all errors but otherwise not responding to them.
func Shutdown() {
	for _, f := range _shutdownFuncs {
		err := f()
		if err != nil {
			log.Print(err)
		}
	}
}

func RegisterDeviceList(typ string, all func() []string, def func() []string) {
	if _devs == nil {
		_devs = make(map[string][]string)
	}
	if _defaults == nil {
		_defaults = make(map[string][]string)
	}
	if _, ok := _devs[typ]; !ok {
		_devs[typ] = []string{}
	}
	_devs[typ] = append(_devs[typ], all()...)
	if _, ok := _defaults[typ]; !ok {
		_defaults[typ] = []string{}
	}
	_defaults[typ] = append(_defaults[typ], def()...)
}

// Return a list of devices registered under domain, where `domain` is one of the
// defined constants in `devices`, e.g., devices.Temperatures.  The
// `enabledOnly` flag determines whether all devices are returned (false), or
// only the ones that have been enabled for the domain.
func Devices(domain string, all bool) []string {
	if all {
		return _devs[domain]
	} else {
		return _defaults[domain]
	}
}
