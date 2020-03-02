// +build windows

package logging

import (
	"os"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func stderrToLogfile(logfile *os.File) {
	// https://groups.google.com/d/msg/golang-nuts/fG8hEAs7ZXs/tahEOuCEPn0J.
	syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(logfile.Fd()), 2, 0)
}
