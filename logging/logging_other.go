// +build !arm64

package logging

import "syscall"
import "os"

func stderrToLogfile(logfile *os.File) {
	syscall.Dup2(int(logfile.Fd()), 2)
}
