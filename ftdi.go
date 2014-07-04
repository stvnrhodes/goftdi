package ftdi

import (
	"errors"
	"fmt"
	"io"
)

/*
#include <ftdi.h>
#include <libusb.h>
*/
import "C"

// Version represents a version of libftdi.
type Version struct {
	Major, Minor, Micro int
	Version             string // Version as (static) string
	Snapshot            string //Git snapshot version if known. Otherwise "unknown" or empty string.
}

// GetVersion gets the current libftdi version.
func GetVersion() Version {
	info := C.ftdi_get_library_version()
	return Version{
		Major:    int(info.major),
		Minor:    int(info.minor),
		Micro:    int(info.micro),
		Version:  C.GoString(info.version_str),
		Snapshot: C.GoString(info.snapshot_str),
	}
}

func getError(ctx *C.struct_ftdi_context) string {
	return C.GoString(C.ftdi_get_error_string(ctx))
}

// Config holds the vendor and product ids for the usb device.
type Config struct {
	// Vendor and Product are IDs uniquely assigned to each USB manufacturer
	Vendor, Product int
	// Baud specifies the baud rate for the serial connection
	Baud int
}

type conn struct {
	ctx *C.struct_ftdi_context
}

// Open returns a serial connection that can be read or written
func Open(c Config) (io.ReadWriteCloser, error) {
	ctx := C.ftdi_new()
	if ctx == nil {
		return nil, errors.New("Could not make new ftdi context")
	}
	C.ftdi_usb_open(ctx, C.int(c.Vendor), C.int(c.Product))

	// TODO(stvn): Support listing all devices, choosing which device to use
	if ret := C.ftdi_usb_open(ctx, C.int(c.Vendor), C.int(c.Product)); ret < 0 {
		defer C.ftdi_free(ctx)
		return nil, fmt.Errorf("unable to open ftdi device: %d (%s)", ret, getError(ctx))
	}

	if f := C.ftdi_set_baudrate(ctx, C.int(c.Baud)); f < 0 {
		defer C.ftdi_free(ctx)
		return nil, fmt.Errorf("unable to set baudrate: %d (%s)", f, getError(ctx))
	}

	if f := C.ftdi_set_line_property(ctx, 8, C.STOP_BIT_1, C.NONE); f < 0 {
		defer C.ftdi_free(ctx)
		return nil, fmt.Errorf("unable to set line parameters: %d (%s)", f, getError(ctx))
	}

	return &conn{ctx: ctx}, nil
}

func libusbErr(n C.int) error {
	switch {
	case n == -666:
		return errors.New("usb device unavailable")
	case n < 0:
		return errors.New(C.GoString(C.libusb_error_name(n)))
	default:
		return nil
	}
}

func (c *conn) Read(p []byte) (int, error) {
	n := C.ftdi_read_data(c.ctx, (*C.uchar)(&p[0]), C.int(len(p)))
	return int(n), libusbErr(n)
}
func (c *conn) Write(p []byte) (int, error) {
	n := C.ftdi_write_data(c.ctx, (*C.uchar)(&p[0]), C.int(len(p)))
	return int(n), libusbErr(n)
}
func (c *conn) Close() error {
	defer C.ftdi_free(c.ctx)
	if ret := C.ftdi_usb_close(c.ctx); ret < 0 {
		return fmt.Errorf("unable to safely close ftdi device: %d (%s)", ret, getError(c.ctx))
	}
	return nil
}
