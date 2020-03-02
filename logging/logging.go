package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/xxxserxxx/gotop/v3"
)

func New(c gotop.Config) (io.WriteCloser, error) {
	// create the log directory
	if err := os.MkdirAll(c.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to make the log directory: %v", err)
	}
	w := &RotateWriter{
		filename:   filepath.Join(c.LogDir, c.LogFile),
		maxLogSize: c.MaxLogSize,
	}
	err := w.rotate()
	if err != nil {
		return nil, err
	}
	// log time, filename, and line number
	log.SetFlags(log.Ltime | log.Lshortfile)
	// log to file
	log.SetOutput(w)

	stderrToLogfile(w.fp)
	return w, nil
}

type RotateWriter struct {
	lock       sync.Mutex
	filename   string // should be set to the actual filename
	fp         *os.File
	maxLogSize int64
}

func (w *RotateWriter) Close() error {
	return w.fp.Close()
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	// Rotate if the log hits the size limit
	s, err := os.Stat(w.filename)
	if err == nil {
		if s.Size() > w.maxLogSize {
			w.rotate()
		}
	}
	return w.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// This will keep three logs
	for i := 1; i > -1; i-- {
		from := fmt.Sprintf("%s.%d", w.filename, i)
		to := fmt.Sprintf("%s.%d", w.filename, i+1)
		// Rename dest file if it already exists
		_, err = os.Stat(from)
		if err == nil {
			err = os.Rename(from, to)
			if err != nil {
				return
			}
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		err = os.Rename(w.filename, fmt.Sprintf("%s.%d", w.filename, 0))
		if err != nil {
			return
		}
	}

	// open the log file
	w.fp, err = os.OpenFile(w.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %v", w.filename, err)
	}

	return nil
}
