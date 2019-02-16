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
	"bytes"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

var errValueNotFound = fmt.Errorf("Value not found")

var sensorW = [4]int32{
	2,  // SENSOR_VOLTS_DC (uV)
	5,  // SENSOR_WATTS (uW)
	7,  // SENSOR_WATTHOUR (uWh)
	10, // SENSOR_INTEGER
}

const (
	sensorA  = 6 // SENSOR_AMPS (uA)
	sensorAH = 8 // SENSOR_AMPHOUR (uAh)
)

type sensordev struct {
	num           int32
	xname         [16]byte
	maxnumt       [21]int32
	sensors_count int32
}

type sensorStatus int32

const (
	unspecified sensorStatus = iota
	ok
	warning
	critical
	unknown
)

type sensor struct {
	desc   [32]byte
	tv     [16]byte // struct timeval
	value  int64
	typ    [4]byte // enum sensor_type
	status sensorStatus
	numt   int32
	flags  int32
}

type interValue struct {
	v *float64
	s *State
	e *error
}

func sysctl(mib []int32, out unsafe.Pointer, n uintptr) syscall.Errno {
	_, _, e := unix.Syscall6(
		unix.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		uintptr(out),
		uintptr(unsafe.Pointer(&n)),
		uintptr(unsafe.Pointer(nil)),
		0,
	)
	return e
}

func readValue(s sensor, div float64) (float64, error) {
	if s.status == unknown {
		return 0, fmt.Errorf("Unknown value received")
	}

	return float64(s.value) / div, nil
}

func readValues(mib []int32, c int32, values map[string]*interValue) {
	var s sensor
	var i int32
	for i = 0; i < c; i++ {
		mib[4] = i

		if err := sysctl(mib, unsafe.Pointer(&s), unsafe.Sizeof(s)); err != 0 {
			for _, value := range values {
				if *value.e == errValueNotFound {
					*value.e = err
				}
			}
		}

		desc := string(s.desc[:bytes.IndexByte(s.desc[:], 0)])
		isState := strings.HasPrefix(desc, "battery ")

		var value *interValue
		var ok bool

		if isState {
			value, ok = values["state"]
		} else {
			value, ok = values[desc]
		}
		if !ok {
			continue
		}

		if isState {
			//TODO:battery idle(?)
			if desc == "battery critical" {
				*value.s, *value.e = Empty, nil
			} else {
				*value.s, *value.e = newState(desc[8:])
			}
			continue
		}

		if strings.HasSuffix(desc, "voltage") {
			*value.v, *value.e = readValue(s, 1000000)
		} else {
			*value.v, *value.e = readValue(s, 1000)
		}
	}
}

func sensordevIter(cb func(sd sensordev, i int, err error) bool) {
	mib := []int32{6, 11, 0}
	var sd sensordev
	var idx int
	var i int32
	for i = 0; ; i++ {
		mib[2] = i

		e := sysctl(mib, unsafe.Pointer(&sd), unsafe.Sizeof(sd))
		if e != 0 {
			if e == unix.ENXIO {
				continue
			}
			if e == unix.ENOENT {
				break
			}
		}

		if bytes.HasPrefix(sd.xname[:], []byte("acpibat")) {
			var err error
			if e != 0 {
				err = e
			}
			if cb(sd, idx, err) {
				return
			}
			idx++
		}
	}
}

func getBattery(sd sensordev) (*Battery, error) {
	b := &Battery{}
	e := ErrPartial{
		Design:        errValueNotFound,
		Full:          errValueNotFound,
		Current:       errValueNotFound,
		ChargeRate:    errValueNotFound,
		State:         errValueNotFound,
		Voltage:       errValueNotFound,
		DesignVoltage: errValueNotFound,
	}

	mib := []int32{6, 11, sd.num, 0, 0}
	for _, w := range sensorW {
		mib[3] = w

		readValues(mib, sd.maxnumt[w], map[string]*interValue{
			"rate":               {v: &b.ChargeRate, e: &e.ChargeRate},
			"design capacity":    {v: &b.Design, e: &e.Design},
			"last full capacity": {v: &b.Full, e: &e.Full},
			"remaining capacity": {v: &b.Current, e: &e.Current},
			"current voltage":    {v: &b.Voltage, e: &e.Voltage},
			"voltage":            {v: &b.DesignVoltage, e: &e.DesignVoltage},
			"state":              {s: &b.State, e: &e.State},
		})
	}

	if e.DesignVoltage != nil && e.Voltage == nil {
		b.DesignVoltage, e.DesignVoltage = b.Voltage, nil
	}

	if e.ChargeRate == errValueNotFound {
		if e.Voltage == nil {
			mib[3] = sensorA

			readValues(mib, sd.maxnumt[sensorA], map[string]*interValue{
				"rate": {v: &b.ChargeRate, e: &e.ChargeRate},
			})

			b.ChargeRate *= b.Voltage
		} else {
			e.ChargeRate = e.Voltage
		}
	}
	if e.Design == errValueNotFound || e.Full == errValueNotFound || e.Current == errValueNotFound {
		mib[3] = sensorAH

		readValues(mib, sd.maxnumt[sensorAH], map[string]*interValue{
			"design capacity":    {v: &b.Design, e: &e.Design},
			"last full capacity": {v: &b.Full, e: &e.Full},
			"remaining capacity": {v: &b.Current, e: &e.Current},
		})

		b.Design *= b.DesignVoltage
		b.Full *= b.Voltage
		b.Current *= b.Voltage
	}

	return b, e
}

func systemGet(idx int) (*Battery, error) {
	var b *Battery
	var e error

	sensordevIter(func(sd sensordev, i int, err error) bool {
		if i == idx {
			if err == nil {
				b, e = getBattery(sd)
			} else {
				e = err
			}
			return true
		}
		return false
	})

	if b == nil {
		return nil, ErrNotFound
	}
	return b, e
}

func systemGetAll() ([]*Battery, error) {
	var batteries []*Battery
	var errors Errors

	sensordevIter(func(sd sensordev, i int, err error) bool {
		var b *Battery
		if err == nil {
			b, err = getBattery(sd)
		}

		batteries = append(batteries, b)
		errors = append(errors, err)
		return false
	})

	return batteries, errors
}
