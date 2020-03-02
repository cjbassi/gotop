package devices

import "log"

var shutdownFuncs []func() error

// RegisterShutdown stores a function to be called by gotop on exit, allowing
// extensions to properly release resources.  Extensions should register a
// shutdown function IFF the extension is using resources that need to be
// released.  The returned error will be logged, but no other action will be
// taken.
func RegisterShutdown(f func() error) {
	shutdownFuncs = append(shutdownFuncs, f)
}

// Shutdown will be called by the `main()` function if gotop is exited
// cleanly.  It will call all of the registered shutdown functions of devices,
// logging all errors but otherwise not responding to them.
func Shutdown() {
	for _, f := range shutdownFuncs {
		err := f()
		if err != nil {
			log.Print(err)
		}
	}
}
