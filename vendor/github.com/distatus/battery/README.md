battery [![Build Status](https://travis-ci.org/distatus/battery.svg?branch=master)](https://travis-ci.org/distatus/battery) [![Go Report Card](https://goreportcard.com/badge/github.com/distatus/battery)](https://goreportcard.com/report/github.com/distatus/battery) [![GoDoc](https://godoc.org/github.com/distatus/battery?status.svg)](https://godoc.org/github.com/distatus/battery)
=======

Cross-platform, normalized battery information library.

Gives access to a system independent, typed battery state, capacity, charge and voltage values recalculated as necessary to be returned in mW, mWh or V units.

Currently supported systems:

* Linux 2.6.39+
* OS X 10.10+
* Windows XP+
* FreeBSD
* DragonFlyBSD
* NetBSD
* OpenBSD
* Solaris

Installation
------------

```bash
$ go get -u github.com/distatus/battery
```

Code Example
------------

```go
import (
	"fmt"

	"github.com/distatus/battery"
)

func main() {
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}
	for i, battery := range batteries {
		fmt.Printf("Bat%d: ", i)
		fmt.Printf("state: %f, ", battery.State)
		fmt.Printf("current capacity: %f mWh, ", battery.Current)
		fmt.Printf("last full capacity: %f mWh, ", battery.Full)
		fmt.Printf("design capacity: %f mWh, ", battery.Design)
		fmt.Printf("charge rate: %f mW, ", battery.ChargeRate)
		fmt.Printf("voltage: %f V, ", battery.Voltage)
		fmt.Printf("design voltage: %f V\n", battery.DesignVoltage)
	}
}
```

CLI
---

There is also a little utility which - more or less - mimicks the GNU/Linux `acpi -b` command.

*Installation*

```bash
$ go get -u github.com/distatus/battery/cmd/battery
```

*Usage*

```bash
$ battery
BAT0: Full, 95.61% [Voltage: 12.15V (design: 12.15V)]
```
