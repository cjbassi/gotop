// +build !darwin,!windows

package appdir

import (
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) UserConfig() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".config")
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		baseDir = os.Getenv("XDG_CONFIG_HOME")
	}

	return filepath.Join(baseDir, d.name)
}

func (d *dirs) UserCache() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".cache")
	if os.Getenv("XDG_CACHE_HOME") != "" {
		baseDir = os.Getenv("XDG_CACHE_HOME")
	}

	return filepath.Join(baseDir, d.name)
}

func (d *dirs) UserLogs() string {
	baseDir := filepath.Join(os.Getenv("HOME"), ".local", "state")
	if os.Getenv("XDG_STATE_HOME") != "" {
		baseDir = os.Getenv("XDG_STATE_HOME")
	}

	return filepath.Join(baseDir, d.name)
}
