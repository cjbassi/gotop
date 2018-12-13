// +build !arm64

package logging

import (
	"os"
	"syscall"
)

func StderrToLogfile(lf *os.File) {
	syscall.Dup2(int(lf.Fd()), 2)
}
