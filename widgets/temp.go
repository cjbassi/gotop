package widgets

// Temp is too customized to inherit from a generic widget so we create a customized one here.
// Temp defines its own Buffer method directly.

import (
	"fmt"
	"sort"
	"strings"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	psHost "github.com/shirou/gopsutil/host"
)

type Temp struct {
	*ui.Block
	interval  time.Duration
	Data      map[string]int
	Threshold int
	TempLow   ui.Color
	TempHigh  ui.Color
}

func NewTemp() *Temp {
	t := &Temp{
		Block:     ui.NewBlock(),
		interval:  time.Second * 5,
		Data:      make(map[string]int),
		Threshold: 80, // temp at which color should change
	}
	t.Label = "Temperatures"

	go t.update()
	ticker := time.NewTicker(t.interval)
	go func() {
		for range ticker.C {
			t.update()
		}
	}()

	return t
}

func (t *Temp) update() {
	sensors, _ := psHost.SensorsTemperatures()
	for _, sensor := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(sensor.SensorKey, "input") {
			// removes '_input' from the end of the sensor name
			label := sensor.SensorKey[:strings.Index(sensor.SensorKey, "_input")]
			t.Data[label] = int(sensor.Temperature)
		}
	}
}

// Buffer implements ui.Bufferer interface.
func (t *Temp) Buffer() *ui.Buffer {
	buf := t.Block.Buffer()

	var keys []string
	for k := range t.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for y, key := range keys {
		if y+1 > t.Y {
			break
		}

		fg := t.TempLow
		if t.Data[key] >= t.Threshold {
			fg = t.TempHigh
		}

		s := ui.MaxString(key, (t.X - 4))
		buf.SetString(1, y+1, s, t.Fg, t.Bg)
		buf.SetString(t.X-2, y+1, fmt.Sprintf("%dC", t.Data[key]), fg, t.Bg)

	}

	return buf
}
