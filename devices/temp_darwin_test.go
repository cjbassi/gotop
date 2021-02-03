// +build darwin

package devices

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func Test_NumCols(t *testing.T) {
	parser := csv.NewReader(bytes.NewReader(smcData))
	parser.Comma = '\t'
	var line []string
	for {
		if line, err = parser.Read(); err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error parsing SMC tags for temp widget: %s", err)
			break
		}
		// The line is malformed if len(line) != 2, but because the asset is static
		// it makes no sense to report the error to downstream users. This must be
		// tested at/around compile time.
		if len(line) == 2 {
			t.Errorf("smc CSV data malformed: expected 2 columns, got %d", len(line))
		}
	}
}

func Test_loadIDs(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"TCAD", "CPU 1 Package Alt."},
		{"TC1P", "CPU 2 Proximity"},
		{"TC1H", "CPU 2 Heatsink"},
		{"TC1D", "CPU 2 Package"},
		{"TC1E", "CPU 2"},
		{"TC1F", "CPU 2"},
		{"TCBH", "CPU 2 Heatsink Alt."},
		{"TCBD", "CPU 2 Package Alt."},
		{"TG0P", "GPU Proximity"},
		{"TG1D", "GPU Die"},
		{"TG0H", "GPU Heatsink"},
		{"TG1H", "GPU Heatsink"},
		{"Ts0S", "Memory Proximity"},
		{"TM0P", "Mem Bank A1"},
		{"TM9P", "Mem Bank B2"},
		{"TCXC", "PECI CPU"},
		{"PSTR", "System Total"},
	}
	ids := loadIDs()
	L := 161
	if len(ids) != L {
		t.Errorf("len(loadIDs) = %d, want %d", len(ids), L)
	}
	for _, tt := range tests {
		t.Run("contents", func(t *testing.T) {
			got := ids[tt.key]
			if got != tt.want {
				t.Errorf("ids[%s] = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
