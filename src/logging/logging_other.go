// +build !arm64

package logging

import (
	"os"
	"syscall"
)

func StderrToLogfile(logfile *os.File) {
	syscall.Dup2(int(logfile.Fd()), 2)
}
