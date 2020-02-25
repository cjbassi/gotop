package devices

var shutdownFuncs []func()

func Shutdown() {
	for _, f := range shutdownFuncs {
		f()
	}
}
