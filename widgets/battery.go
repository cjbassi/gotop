package widgets

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/distatus/battery"

	ui "github.com/xxxserxxx/gotop/v4/termui"
)

type BatteryWidget struct {
	*ui.LineGraph
	updateInterval time.Duration
}

func NewBatteryWidget(horizontalScale int) *BatteryWidget {
	self := &BatteryWidget{
		LineGraph:      ui.NewLineGraph(),
		updateInterval: time.Minute,
	}
	self.Title = tr.Value("widget.label.battery")
	self.HorizontalScale = horizontalScale

	// intentional duplicate
	// adds 2 datapoints to the graph, otherwise the dot is difficult to see
	self.update()
	self.update()

	go func() {
		for range time.NewTicker(self.updateInterval).C {
			self.Lock()
			self.update()
			self.Unlock()
		}
	}()

	return self
}

func (b *BatteryWidget) EnableMetric() {
	bats, err := battery.GetAll()
	if err != nil {
		log.Printf(tr.Value("error.metricsetup", "batt", err.Error()))
		return
	}
	for i, _ := range bats {
		id := makeID(i)
		metrics.NewGauge(makeName("battery", i), func() float64 {
			if ds, ok := b.Data[id]; ok {
				return ds[len(ds)-1]
			}
			return 0.0
		})
	}
}

func makeID(i int) string {
	return tr.Value("widget.label.batt") + strconv.Itoa(i)
}

func (b *BatteryWidget) Scale(i int) {
	b.LineGraph.HorizontalScale = i
}

func (b *BatteryWidget) update() {
	batteries, err := battery.GetAll()
	if err != nil {
		switch errt := err.(type) {
		case battery.ErrFatal:
			log.Printf(tr.Value("error.fatalfetch", "batt", err.Error()))
			return
		case battery.Errors:
			batts := make([]*battery.Battery, 0)
			for i, e := range errt {
				if e == nil {
					batts = append(batts, batteries[i])
				} else {
					log.Printf(tr.Value("error.recovfetch"), "batt", e.Error())
				}
			}
			if len(batts) < 1 {
				log.Print(tr.Value("error.nodevfound", "batt"))
				return
			}
			batteries = batts
		}
	}
	for i, battery := range batteries {
		id := makeID(i)
		perc := battery.Current / battery.Full
		percentFull := math.Abs(perc) * 100.0
		// TODO: look into this sort of thing; doesn't the array grow forever? Is the widget library truncating it?
		b.Data[id] = append(b.Data[id], percentFull)
		b.Labels[id] = fmt.Sprintf("%3.0f%% %.0f/%.0f", percentFull, math.Abs(battery.Current), math.Abs(battery.Full))
	}
}
