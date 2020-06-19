package devices

import (
	"log"
)

// TODO add thermal history graph. Update when something changes?

var tempUpdates []func(map[string]int) map[string]error

func RegisterTemp(update func(map[string]int) map[string]error) {
	tempUpdates = append(tempUpdates, update)
}

func UpdateTemps(temps map[string]int) {
	for _, f := range tempUpdates {
		errs := f(temps)
		if errs != nil {
			for k, e := range errs {
				log.Printf("error updating temp for %s: %s", k, e)
			}
		}
	}
}
