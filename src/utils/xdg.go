package utils

import (
	"os"
	"path/filepath"
)

func GetConfigDir(name string) string {
	var basedir string
	if env := os.Getenv("XDG_CONFIG_HOME"); env != "" {
		basedir = env
	} else {
		basedir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(basedir, name)
}

func GetLogDir(name string) string {
	var basedir string
	if env := os.Getenv("XDG_STATE_HOME"); env != "" {
		basedir = env
	} else {
		basedir = filepath.Join(os.Getenv("HOME"), ".local", "state")
	}
	return filepath.Join(basedir, name)
}
