package widgets

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/distatus/battery"

	ui "github.com/xxxserxxx/gotop/termui"
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
	self.Title = " Battery Status "
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

func makeId(i int) string {
	return "Batt" + strconv.Itoa(i)
}

func (b *BatteryWidget) Scale(i int) {
	b.LineGraph.HorizontalScale = i
}

func (self *BatteryWidget) update() {
	batteries, err := battery.GetAll()
	if err != nil {
		switch errt := err.(type) {
		case battery.ErrFatal:
			log.Printf("fatal error fetching battery info: %v", err)
			return
		case battery.Errors:
			batts := make([]*battery.Battery, 0)
			for i, e := range errt {
				if e == nil {
					batts = append(batts, batteries[i])
				} else {
					log.Printf("recoverable error fetching battery info; skipping battery: %v", e)
				}
			}
			if len(batts) < 1 {
				log.Print("no usable batteries found")
				return
			}
			batteries = batts
		}
	}
	for i, battery := range batteries {
		id := makeId(i)
		percentFull := math.Abs(battery.Current/battery.Full) * 100.0
		self.Data[id] = append(self.Data[id], percentFull)
		self.Labels[id] = fmt.Sprintf("%3.0f%% %.0f/%.0f", percentFull, math.Abs(battery.Current), math.Abs(battery.Full))
	}
}
