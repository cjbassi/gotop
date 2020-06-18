package colorschemes

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shibukawa/configdir"
)

var registry map[string]Colorscheme

func init() {
	if registry == nil {
		registry = make(map[string]Colorscheme)
	}
}

// FromName loads a Colorscheme by name; confDir is used to search
// directories for a scheme matching the name.  The search order
// is the same as for config files.
func FromName(confDir configdir.ConfigDir, c string) (Colorscheme, error) {
	if cs, ok := registry[c]; ok {
		return cs, nil
	}
	cs, err := getCustomColorscheme(confDir, c)
	return cs, err
}

func register(name string, c Colorscheme) {
	if registry == nil {
		registry = make(map[string]Colorscheme)
	}
	c.Name = name
	registry[name] = c
}

// getCustomColorscheme	tries to read a custom json colorscheme from <configDir>/<name>.json
func getCustomColorscheme(confDir configdir.ConfigDir, name string) (Colorscheme, error) {
	var cs Colorscheme
	fn := name + ".json"
	folder := confDir.QueryFolderContainsFile(fn)
	if folder == nil {
		paths := make([]string, 0)
		for _, d := range confDir.QueryFolders(configdir.Existing) {
			paths = append(paths, d.Path)
		}
		return cs, fmt.Errorf("failed to find colorscheme file %s in %s", fn, strings.Join(paths, ", "))
	}
	dat, err := folder.ReadFile(fn)
	if err != nil {
		return cs, fmt.Errorf("failed to read colorscheme file %s: %v", filepath.Join(folder.Path, fn), err)
	}
	err = json.Unmarshal(dat, &cs)
	if err != nil {
		return cs, fmt.Errorf("failed to parse colorscheme file: %v", err)
	}
	return cs, nil
}
