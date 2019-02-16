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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os/exec"
	"strconv"
)

var errValueNotFound = fmt.Errorf("Value not found")

func readFloat(val string) (float64, error) {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	if num == math.MaxUint32 {
		return 0, fmt.Errorf("Unknown value received")
	}
	return num, nil
}

func readVoltage(val string) (float64, error) {
	voltage, err := readFloat(val)
	if err != nil {
		return 0, err
	}
	return voltage / 1000, nil
}

func readState(val string) (State, error) {
	state, err := strconv.Atoi(val)
	if err != nil {
		return Unknown, err
	}

	switch {
	case state&1 != 0:
		return newState("Discharging")
	case state&2 != 0:
		return newState("Charging")
	case state&4 != 0:
		return newState("Empty")
	default:
		return Unknown, fmt.Errorf("Invalid state flag retrieved: `%d`", state)
	}
}

type errParse int

func (p errParse) Error() string {
	return fmt.Sprintf("Parse error: `%d`", p)
}

type batteryReader struct {
	cmdout *bufio.Scanner
	li     int
	lline  []byte
	e      ErrPartial
}

func (r *batteryReader) setErrParse(n int) {
	if r.e.Design == errValueNotFound {
		r.e.Design = errParse(n)
	}
	if r.e.Full == errValueNotFound {
		r.e.Full = errParse(n)
	}
	if r.e.Current == errValueNotFound {
		r.e.Current = errParse(n)
	}
	if r.e.ChargeRate == errValueNotFound {
		r.e.ChargeRate = errParse(n)
	}
	if r.e.State == errValueNotFound {
		r.e.State = errParse(n)
	}
	if r.e.Voltage == errValueNotFound {
		r.e.Voltage = errParse(n)
	}
	if r.e.DesignVoltage == errValueNotFound {
		r.e.DesignVoltage = errParse(n)
	}
}

func (r *batteryReader) readValue() (string, string, int) {
	var piece []byte
	if r.lline != nil {
		piece = r.lline
		r.lline = nil
	} else {
		pieces := bytes.Split(r.cmdout.Bytes(), []byte{':'})
		if len(pieces) < 4 {
			return "", "", 4
		}

		i, err := strconv.Atoi(string(pieces[1]))
		if err != nil {
			return "", "", 1
		}

		if i != r.li {
			r.li = i
			r.lline = pieces[3]
			return "", "", 666
		}

		piece = pieces[3]
	}

	values := bytes.Split(piece, []byte{'\t'})
	if len(values) < 2 {
		return "", "", 2
	}
	return string(values[0]), string(values[1]), 0
}

func (r *batteryReader) readBattery() (*Battery, bool, bool) {
	b := &Battery{}
	var exists, amps bool

	for r.cmdout.Scan() {
		exists = true

		name, value, errno := r.readValue()
		if errno == 666 {
			break
		}
		if errno != 0 {
			r.setErrParse(errno)
			continue
		}

		switch name {
		case "bif_design_cap":
			b.Design, r.e.Design = readFloat(value)
		case "bif_last_cap":
			b.Full, r.e.Full = readFloat(value)
		case "bif_unit":
			amps = value != "0"
		case "bif_voltage":
			b.DesignVoltage, r.e.DesignVoltage = readVoltage(value)
		case "bst_voltage":
			b.Voltage, r.e.Voltage = readVoltage(value)
		case "bst_rem_cap":
			b.Current, r.e.Current = readFloat(value)
		case "bst_rate":
			b.ChargeRate, r.e.ChargeRate = readFloat(value)
		case "bst_state":
			b.State, r.e.State = readState(value)
		}
	}

	return b, amps, exists
}

func (r *batteryReader) next() (*Battery, error) {
	r.e = ErrPartial{
		Design:        errValueNotFound,
		Full:          errValueNotFound,
		Current:       errValueNotFound,
		ChargeRate:    errValueNotFound,
		State:         errValueNotFound,
		Voltage:       errValueNotFound,
		DesignVoltage: errValueNotFound,
	}

	b, amps, exists := r.readBattery()

	if !exists {
		return nil, io.EOF
	}

	if r.e.DesignVoltage != nil && r.e.Voltage == nil {
		b.DesignVoltage, r.e.DesignVoltage = b.Voltage, nil
	}

	if amps {
		if r.e.DesignVoltage == nil {
			b.Design *= b.DesignVoltage
		} else {
			r.e.Design = r.e.DesignVoltage
		}
		if r.e.Voltage == nil {
			b.Full *= b.Voltage
			b.Current *= b.Voltage
			b.ChargeRate *= b.Voltage
		} else {
			r.e.Full = r.e.Voltage
			r.e.Current = r.e.Voltage
			r.e.ChargeRate = r.e.Voltage
		}
	}

	if b.State == Unknown && r.e.Current == nil && r.e.Full == nil && b.Current >= b.Full {
		b.State, r.e.State = newState("Full")
	}

	return b, r.e
}

func newBatteryReader() (*batteryReader, error) {
	out, err := exec.Command("kstat", "-p", "-m", "acpi_drv", "-n", "battery B*").Output()
	if err != nil {
		return nil, err
	}

	return &batteryReader{cmdout: bufio.NewScanner(bytes.NewReader(out))}, nil
}

func systemGet(idx int) (*Battery, error) {
	br, err := newBatteryReader()
	if err != nil {
		return nil, err
	}

	return br.next()
}

func systemGetAll() ([]*Battery, error) {
	br, err := newBatteryReader()
	if err != nil {
		return nil, err
	}

	var batteries []*Battery
	var errors Errors
	for b, e := br.next(); e != io.EOF; b, e = br.next() {
		batteries = append(batteries, b)
		errors = append(errors, e)
	}

	return batteries, errors
}
