// +build darwin

package devices

import (
	"testing"
)

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
