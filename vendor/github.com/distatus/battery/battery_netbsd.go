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
	"errors"
	"math"
	"sort"
	"strings"
	"unsafe"

	plist "howett.net/plist"

	"golang.org/x/sys/unix"
)

type plistref struct {
	pref_plist unsafe.Pointer
	pref_len   uint64
}

type values struct {
	Description string `plist:"description"`
	CurValue    int    `plist:"cur-value"`
	MaxValue    int    `plist:"max-value"`
	State       string `plist:"state"`
	Type        string `plist:"type"`
}

type prop []values

type props map[string]prop

func readBytes(ptr unsafe.Pointer, length uint64) []byte {
	buf := make([]byte, length-1)
	var i uint64
	for ; i < length-1; i++ {
		buf[i] = *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)))
	}
	return buf
}

func readProps() (props, error) {
	fd, err := unix.Open("/dev/sysmon", unix.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer unix.Close(fd)

	var retptr plistref

	if err = ioctl(fd, 0, 'E', unsafe.Sizeof(retptr), unsafe.Pointer(&retptr)); err != nil {
		return nil, err
	}
	bytes := readBytes(retptr.pref_plist, retptr.pref_len)

	var props props
	if _, err = plist.Unmarshal(bytes, &props); err != nil {
		return nil, err
	}
	return props, nil
}

func handleValue(val values, div float64, res *float64, amps *[]string) error {
	if val.State == "invalid" || val.State == "unknown" {
		return errors.New("Unknown value received")
	}

	*res = float64(val.CurValue) / div

	if amps != nil && strings.HasPrefix(val.Type, "Amp") {
		*amps = append(*amps, val.Description)
	}

	return nil
}

func deriveState(cr1, cr2 error, current float64, max int) (State, error) {
	if cr1 == nil && cr2 != nil {
		return Charging, nil
	}
	if cr1 != nil && cr2 == nil {
		return Discharging, nil
	}
	if cr1 != nil && cr2 != nil && current == float64(max)/1000 {
		return Full, nil
	}
	return Unknown, errors.New("Contradicting values received")
}

func handleVoltage(amps []string, b *Battery, e *ErrPartial) {
	if e.DesignVoltage != nil && e.Voltage == nil {
		b.DesignVoltage, e.DesignVoltage = b.Voltage, nil
	}

	for _, val := range amps {
		switch val {
		case "design cap":
			if e.DesignVoltage == nil {
				b.Design *= b.DesignVoltage
			} else {
				e.Design = e.DesignVoltage
			}
		case "last full cap":
			if e.Voltage == nil {
				b.Full *= b.Voltage
			} else {
				e.Full = e.Voltage
			}
		case "charge":
			if e.Voltage == nil {
				b.Current *= b.Voltage
			} else {
				e.Current = e.Voltage
			}
		case "charge rate", "discharge rate":
			if e.Voltage == nil {
				b.ChargeRate *= b.Voltage
			} else {
				e.ChargeRate = e.Voltage
			}
		}
	}
}

func sortFilterProps(props props) []string {
	var keys []string
	for key := range props {
		if key[:7] != "acpibat" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func convertBattery(prop prop) (*Battery, error) {
	b := &Battery{}
	e := ErrPartial{}

	amps := []string{}
	var cr1, cr2 error
	var maxCharge int

	for _, val := range prop {
		switch val.Description {
		case "voltage":
			e.Voltage = handleValue(val, 1000000, &b.Voltage, nil)
		case "design voltage":
			e.DesignVoltage = handleValue(val, 1000000, &b.DesignVoltage, nil)
		case "design cap":
			e.Design = handleValue(val, 1000, &b.Design, &amps)
		case "last full cap":
			e.Full = handleValue(val, 1000, &b.Full, &amps)
		case "charge":
			e.Current = handleValue(val, 1000, &b.Current, &amps)
			maxCharge = val.MaxValue
		case "charge rate":
			cr1 = handleValue(val, 1000, &b.ChargeRate, &amps)
		case "discharge rate":
			cr2 = handleValue(val, 1000, &b.ChargeRate, &amps)
			b.ChargeRate = math.Abs(b.ChargeRate)
		}
	}

	b.State, e.State = deriveState(cr1, cr2, b.Current, maxCharge)

	handleVoltage(amps, b, &e)

	return b, e
}

func systemGet(idx int) (*Battery, error) {
	props, err := readProps()
	if err != nil {
		return nil, err
	}

	keys := sortFilterProps(props)
	if idx >= len(keys) {
		return nil, ErrNotFound
	}
	return convertBattery(props[keys[idx]])
}

func systemGetAll() ([]*Battery, error) {
	props, err := readProps()
	if err != nil {
		return nil, err
	}

	keys := sortFilterProps(props)
	batteries := make([]*Battery, len(keys))
	errors := make(Errors, len(keys))
	for i, key := range keys {
		batteries[i], errors[i] = convertBattery(props[key])
	}

	return batteries, errors
}
