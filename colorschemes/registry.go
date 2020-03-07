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

func FromName(confDir configdir.ConfigDir, c string) (Colorscheme, error) {
	cs, ok := registry[c]
	if !ok {
		cs, err := getCustomColorscheme(confDir, c)
		if err != nil {
			return cs, err
		}
	}
	return cs, nil
}

func register(name string, c Colorscheme) {
	if registry == nil {
		registry = make(map[string]Colorscheme)
	}
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
