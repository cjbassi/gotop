package widgets

// loosely based on https://github.com/openbsd/src/blob/master/sbin/sysctl/sysctl.c#L2517

// #include <sys/time.h>
// #include <sys/sysctl.h>
// #include <sys/sensors.h>
import "C"

import (
	"strconv"
	"syscall"
	"unsafe"

	"github.com/cjbassi/gotop/src/utils"
)

func (self *TempWidget) getTemp(mib []C.int, mlen int, snsrdev *C.struct_sensordev, index int) {
	switch mlen {
	case 4:
		k := mib[3]
		var numt C.int
		for numt = 0; numt < snsrdev.maxnumt[k]; numt++ {
			mib[4] = numt
			self.getTemp(mib, mlen+1, snsrdev, int(numt))
		}
	case 5:
		var snsr C.struct_sensor
		var slen C.size_t = C.sizeof_struct_sensor

		if v, _ := C.sysctl(&mib[0], 5, unsafe.Pointer(&snsr), &slen, nil, 0); v == -1 {
			return
		}

		if slen > 0 && (snsr.flags&C.SENSOR_FINVALID) == 0 {
			key := C.GoString(&snsrdev.xname[0]) + ".temp" + strconv.Itoa(index)
			temp := int((snsr.value - 273150000.0) / 1000000.0)

			switch self.TempScale {
			case Fahrenheit:
				self.Data[key] = utils.CelsiusToFahrenheit(temp)
			case Celcius:
				self.Data[key] = temp
			}
		}
	}
}

func (self *TempWidget) update() {
	mib := []C.int{0, 1, 2, 3, 4}

	var snsrdev C.struct_sensordev
	var len C.ulong = C.sizeof_struct_sensordev

	mib[0] = C.CTL_HW
	mib[1] = C.HW_SENSORS
	mib[3] = C.SENSOR_TEMP

	var i C.int
	for i = 0; ; i++ {
		mib[2] = i
		if v, e := C.sysctl(&mib[0], 3, unsafe.Pointer(&snsrdev), &len, nil, 0); v == -1 {
			if e == syscall.ENXIO {
				continue
			}
			if e == syscall.ENOENT {
				break
			}
		}
		self.getTemp(mib, 4, &snsrdev, 0)
	}
}
