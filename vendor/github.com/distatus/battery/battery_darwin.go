// battery
// Copyright (C) 2016-2017 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package battery

import (
	"math"
	"os/exec"

	plist "howett.net/plist"
)

type battery struct {
	Voltage           int
	CurrentCapacity   int
	MaxCapacity       int
	DesignCapacity    int
	Amperage          int64
	FullyCharged      bool
	IsCharging        bool
	ExternalConnected bool
}

func readBatteries() ([]*battery, error) {
	out, err := exec.Command("ioreg", "-n", "AppleSmartBattery", "-r", "-a").Output()
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		// No batteries.
		return nil, nil
	}

	var data []*battery
	if _, err = plist.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func convertBattery(battery *battery) *Battery {
	volts := float64(battery.Voltage) / 1000
	b := &Battery{
		Current:       float64(battery.CurrentCapacity) * volts,
		Full:          float64(battery.MaxCapacity) * volts,
		Design:        float64(battery.DesignCapacity) * volts,
		ChargeRate:    math.Abs(float64(battery.Amperage)) * volts,
		Voltage:       volts,
		DesignVoltage: volts,
	}
	switch {
	case !battery.ExternalConnected:
		b.State, _ = newState("Discharging")
	case battery.IsCharging:
		b.State, _ = newState("Charging")
	case battery.CurrentCapacity == 0:
		b.State, _ = newState("Empty")
	case battery.FullyCharged:
		b.State, _ = newState("Full")
	default:
		b.State, _ = newState("Unknown")
	}
	return b
}

func systemGet(idx int) (*Battery, error) {
	batteries, err := readBatteries()
	if err != nil {
		return nil, err
	}

	if idx >= len(batteries) {
		return nil, ErrNotFound
	}
	return convertBattery(batteries[idx]), nil
}

func systemGetAll() ([]*Battery, error) {
	_batteries, err := readBatteries()
	if err != nil {
		return nil, err
	}

	batteries := make([]*Battery, len(_batteries))
	for i, battery := range _batteries {
		batteries[i] = convertBattery(battery)
	}
	return batteries, nil
}
