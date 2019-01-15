package appdir

import (
	"os"
	"path/filepath"
)

type dirs struct {
	name string
}

func (d *dirs) UserConfig() string {
	return filepath.Join(os.Getenv("APPDATA"), d.name)
}

func (d *dirs) UserCache() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), d.name)
}

func (d *dirs) UserLogs() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), d.name)
}
