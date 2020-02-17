package colorschemes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var registry map[string]Colorscheme

func FromName(confDir string, c string) (Colorscheme, error) {
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
func getCustomColorscheme(confDir string, name string) (Colorscheme, error) {
	var cs Colorscheme
	filePath := filepath.Join(confDir, name+".json")
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return cs, fmt.Errorf("failed to read colorscheme file: %v", err)
	}
	err = json.Unmarshal(dat, &cs)
	if err != nil {
		return cs, fmt.Errorf("failed to parse colorscheme file: %v", err)
	}
	return cs, nil
}
