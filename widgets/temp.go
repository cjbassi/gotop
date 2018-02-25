package widgets

// Temp is too customized to inherit from a generic widget so we create a customized one here.
// Temp defines its own Buffer method directly.

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/cjbassi/gotop/termui"
	ps "github.com/shirou/gopsutil/host"
)

type Temp struct {
	*ui.Block
	interval   time.Duration
	Data       []int
	DataLabels []string
	Threshold  int
	TempLow    ui.Color
	TempHigh   ui.Color
}

func NewTemp() *Temp {
	t := &Temp{
		Block:     ui.NewBlock(),
		interval:  time.Second * 5,
		Threshold: 80, // temp at which color should change to red
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
	sensors, _ := ps.SensorsTemperatures()
	temps := []int{}
	labels := []string{}
	for _, temp := range sensors {
		// only sensors with input in their name are giving us live temp info
		if strings.Contains(temp.SensorKey, "input") {
			temps = append(temps, int(temp.Temperature))
			// removes '_input' from the end of the sensor name
			labels = append(labels, temp.SensorKey[:strings.Index(temp.SensorKey, "_input")])
		}
	}
	t.Data = temps
	t.DataLabels = labels
}

// Buffer implements ui.Bufferer interface.
func (t *Temp) Buffer() *ui.Buffer {
	buf := t.Block.Buffer()

	for y, text := range t.DataLabels {
		if y+1 > t.Y {
			break
		}

		fg := t.TempLow
		if t.Data[y] >= t.Threshold {
			fg = t.TempHigh
		}

		s := ui.MaxString(text, (t.X - 4))
		buf.SetString(1, y+1, s, t.Fg, t.Bg)
		buf.SetString(t.X-2, y+1, fmt.Sprintf("%dC", t.Data[y]), fg, t.Bg)
	}

	return buf
}
