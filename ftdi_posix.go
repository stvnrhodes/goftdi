// +build !windows,cgo

package ftdi

/*
#cgo CFLAGS: -I /usr/include/libftdi1 -I /usr/local/include/libusb-1.0
#cgo LDFLAGS: -lftdi1 -lusb-1.0
*/
import "C"
