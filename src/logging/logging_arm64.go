package logging

import (
	"os"
	"syscall"
)

func StderrToLogfile(lf *os.File) {
	syscall.Dup3(int(lf.Fd()), 2, 0)
}
