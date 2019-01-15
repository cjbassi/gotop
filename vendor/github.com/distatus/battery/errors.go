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

import "fmt"

// ErrNotFound variable represents battery not found error.
//
// Only ever returned wrapped in ErrFatal.
var ErrNotFound = fmt.Errorf("Not found")

// ErrAllNotNil variable says that backend returned ErrPartial with
// all fields having not nil values, hence it was converted to ErrFatal.
//
// Only ever returned wrapped in ErrFatal.
var ErrAllNotNil = fmt.Errorf("All fields had not nil errors")

// ErrFatal type represents a fatal error.
//
// It indicates that either the library was not able to perform some kind
// of operation critical to retrieving any data, or all partials have failed at
// once (which would be equivalent to returning a ErrPartial with no nils).
//
// As such, the caller should assume that no meaningful data was
// returned alongside the error and act accordingly.
type ErrFatal struct {
	Err error // The actual error that happened.
}

func (f ErrFatal) Error() string {
	return fmt.Sprintf("Could not retrieve battery info: `%s`", f.Err)
}

// ErrPartial type represents a partial error.
//
// It indicates that there were problems retrieving some of the data,
// but some was also retrieved successfully.
// If there would be all nils, nil is returned instead.
// If there would be all not nils, ErrFatal is returned instead.
//
// The fields represent fields in the Battery type.
type ErrPartial struct {
	State         error
	Current       error
	Full          error
	Design        error
	ChargeRate    error
	Voltage       error
	DesignVoltage error
}

func (p ErrPartial) Error() string {
	if p.isNil() {
		return "{}"
	}
	errors := map[string]error{
		"State":         p.State,
		"Current":       p.Current,
		"Full":          p.Full,
		"Design":        p.Design,
		"ChargeRate":    p.ChargeRate,
		"Voltage":       p.Voltage,
		"DesignVoltage": p.DesignVoltage,
	}
	keys := []string{"State", "Current", "Full", "Design", "ChargeRate", "Voltage", "DesignVoltage"}
	s := "{"
	for _, name := range keys {
		err := errors[name]
		if err != nil {
			s += fmt.Sprintf("%s:%s ", name, err.Error())
		}
	}
	return s[:len(s)-1] + "}"
}

func (p ErrPartial) isNil() bool {
	return p.State == nil &&
		p.Current == nil &&
		p.Full == nil &&
		p.Design == nil &&
		p.ChargeRate == nil &&
		p.Voltage == nil &&
		p.DesignVoltage == nil
}

func (p ErrPartial) noNil() bool {
	return p.State != nil &&
		p.Current != nil &&
		p.Full != nil &&
		p.Design != nil &&
		p.ChargeRate != nil &&
		p.Voltage != nil &&
		p.DesignVoltage != nil
}

// Errors type represents an array of ErrFatal, ErrPartial or nil values.
//
// Can only possibly be returned by GetAll() call.
type Errors []error

func (e Errors) Error() string {
	s := "["
	for _, err := range e {
		if err != nil {
			s += err.Error() + " "
		}
	}
	if len(s) > 1 {
		s = s[:len(s)-1]
	}
	return s + "]"
}

func wrapError(err error) error {
	if perr, ok := err.(ErrPartial); ok {
		if perr.isNil() {
			return nil
		}
		if perr.noNil() {
			return ErrFatal{ErrAllNotNil}
		}
		return perr
	}
	if err != nil {
		return ErrFatal{err}
	}
	return nil
}
